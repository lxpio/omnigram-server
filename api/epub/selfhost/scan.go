// scan 目录扫描搜刮工具
package selfhost

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nexptr/omnigram-server/api/epub/schema"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nutsdb/nutsdb"
)

// Scanner  文件扫描器
type Scanner struct {
	Running bool     `json:"running"`
	Count   int      `json:"count"` //扫描文件计数
	Errs    []string `json:"errors"`
	root    string   `json:"-"` //扫描错误详情信息
	cached  *nutsdb.DB
	wg      *sync.WaitGroup `json:"-"`
}

func (m *Scanner) Stop() {
	if m.wg != nil {
		m.wg.Wait()
	}
	if m.cached != nil {
		m.cached.Close()
	}
}

func NewScan(root, meta string) (*Scanner, error) {

	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(meta),
	)

	if err != nil {
		log.E("打开metadata失败：", err.Error())
		return nil, err
	}

	return &Scanner{
		Count:  0,
		root:   root,
		cached: db,
		wg:     new(sync.WaitGroup),
		Errs:   []string{},
	}, nil

}

func (m *Scanner) startSingleThread(manager *ScannerManager, books <-chan *schema.Book) {

	errChan := make(chan string)

	m.wg.Add(1)
	go func() {

		defer func() {

			m.wg.Done()
			close(errChan)
			log.I(`退出扫描程序`)
		}()

		for {

			select {

			case <-manager.ctx.Done():
				log.W(`接收到退出命令，退出扫描`)
				return
			case book, ok := <-books:

				if !ok {
					//books is closed
					log.I(`书籍为空，退出解析文件。`)
					return
				}

				log.D(`开始解析: `, book.Path, ` 到数据库`)

				if err := book.GetMetadataFromFile(); err != nil {
					log.E(`获取图书基本元素失败 `, err.Error())
					errChan <- `文件：` + book.Path + ` 解析失败：` + err.Error()
				} else {
					if err := book.Save(manager.ctx, manager.orm, manager.kv); err != nil {
						errChan <- err.Error()

					} else {
						m.cacheEpubFilePath(book.Path)
					}
					//

				}

			}
		}

	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		for err := range errChan {
			m.Errs = append(m.Errs, err)
		}

		//更新扫描状态
		manager.updateStatus(ScanStatus{
			Running:   false,
			BookCount: m.Count,
			Errs:      m.Errs,
		})

		//关闭扫描器
		m.cached.Close()
		m.cached = nil
	}()

}

func (m *Scanner) Start(manager *ScannerManager, maxThread int, refresh bool) {

	books := m.Walk(manager.ctx, refresh)

	if maxThread < 2 {
		m.startSingleThread(manager, books)
		return
	}

	errChan := make(chan string)

	m.wg.Add(1)
	go func() {

		wg := new(sync.WaitGroup)

		defer func() {

			m.wg.Done()

			wg.Wait()

			close(errChan)
			log.I(`退出扫描程序`)
		}()

		concurrent := make(chan struct{}, maxThread)

		for {

			select {

			case <-manager.ctx.Done():
				log.W(`接收到退出命令，退出扫描`)
				return
			case book, ok := <-books:

				if !ok {
					//books is closed
					log.I(`书籍为空，退出解析文件。`)

					return
				}

				log.D(`开始解析: `, book.Path, ` 到数据库`)

				wg.Add(1)
				concurrent <- struct{}{}

				go func(b *schema.Book) {

					defer wg.Done()

					if err := b.GetMetadataFromFile(); err != nil {
						log.E(`获取图书基本元素失败 `, err.Error())
						errChan <- `文件：` + b.Path + ` 解析失败：` + err.Error()
					} else {
						if err := b.Save(manager.ctx, manager.orm, manager.kv); err != nil {
							errChan <- err.Error()
						} else {
							m.cacheEpubFilePath(book.Path)
						}
					}

					<-concurrent

				}(book)

			}
		}

	}()

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		for err := range errChan {
			m.Errs = append(m.Errs, err)
		}

		status := ScanStatus{
			Running:   false,
			BookCount: m.Count,
			Errs:      m.Errs,
		}
		//更新扫描状态
		manager.updateStatus(status)
		m.dumpStats(status)
		m.cached.Close()
		m.cached = nil
	}()

}

// Walk 遍历扫描路径下epub文件
func (m *Scanner) Walk(ctx context.Context, refresh bool) <-chan *schema.Book {

	log.I(`开始扫描路径:`, m.root)
	books := make(chan *schema.Book)

	go func() {

		err := filepath.Walk(m.root, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				log.E(`扫描路径失败：`, path)
				return err
			}

			//只扫描epub文件
			if !info.IsDir() && filepath.Ext(info.Name()) == `.epub` {

				if m.epubFilePathExists(path) && !refresh {
					log.I(`文件：`, path, `已经存在,放弃扫描`)
					return nil
				}

				log.I(`扫描的到文件：`, path)
				book := &schema.Book{
					ID:            0,
					Size:          info.Size(),
					Path:          path,
					CTime:         time.Now(),
					UTime:         time.Now(),
					Rating:        0,
					PublishDate:   `1970-01-01`,
					CountVisit:    0,
					CountDownload: 0,
				}

				books <- book
				m.Count++
			}
			return nil
		})

		if err != nil {
			log.E(`扫描路径失败：`, err.Error())
		}
		close(books)
	}()

	return books
}

func (m *Scanner) cacheEpubFilePath(path string) error {
	return m.cached.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(`epub`, []byte(path), []byte{}, 0); err != nil {
				return err
			}
			return nil
		})

}

func (m *Scanner) epubFilePathExists(path string) bool {

	err := m.cached.View(func(tx *nutsdb.Tx) error {

		_, err := tx.Get(`epub`, []byte(path))
		return err

	})

	return err == nil

}

func (m *Scanner) dumpStats(status ScanStatus) error {

	bytes, _ := json.Marshal(status)

	return m.cached.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(`sys`, []byte("last_scan_status"), bytes, 0); err != nil {
				return err
			}
			return nil
		})

}
