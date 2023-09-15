package selfhost

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/store"
	"github.com/nexptr/omnigram-server/utils"
	"github.com/nutsdb/nutsdb"
	"gorm.io/gorm"
)

// const statsCachePath = `config`

// ScanStatus 扫描状态
type ScanStatus struct {
	Running bool     `json:"running"`
	Count   int      `json:"count"`
	Errs    []string `json:"errs"`
}

type ScannerManager struct {
	// cf *conf.Config
	dataPath string
	kv       store.KV
	orm      *gorm.DB

	ctx context.Context

	sync.RWMutex

	stats ScanStatus
}

func NewScannerManager(ctx context.Context, dataPath string, kv store.KV, orm *gorm.DB) (*ScannerManager, error) {

	//初始化上传文件目录
	uploadPath := filepath.Join(dataPath, `upload`)
	os.MkdirAll(uploadPath, 0755)

	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(filepath.Join(dataPath, utils.ConfigBucket)),
	)

	if err != nil {
		return nil, err
	}

	scanner := &ScannerManager{
		dataPath: dataPath,
		kv:       kv,
		orm:      orm,
		stats:    loadLastScanStatus(db),
		ctx:      ctx,
	}

	//获取本地存储的状态

	return scanner, nil
}

func (m *ScannerManager) IsRunning() bool {
	m.RLock()
	defer m.RUnlock()
	return m.stats.Running
}

func (m *ScannerManager) Status() ScanStatus {
	m.RLock()
	defer m.RUnlock()

	s := ScanStatus{
		Running: m.stats.Running,
		Count:   m.stats.Count,
		Errs:    m.stats.Errs,
	}

	return s

}

func (m *ScannerManager) Start(maxThread int, refresh bool) {

	if m.IsRunning() {
		log.E(`扫描器已经在执行，放弃执行`)
		return
	}
	log.I(`启动文件目录扫描`)
	m.newScan(m.dataPath, maxThread, refresh)

}

func (m *ScannerManager) newScan(path string, maxThread int, refresh bool) {
	m.Lock()
	scan, err := NewScan(path) //new scanner

	if err != nil {
		m.Unlock()
		log.E(err.Error())
		return
	}
	m.stats.Running = true
	m.Unlock()

	scan.Start(m, maxThread, refresh)
}

func (m *ScannerManager) updateStatus(stats ScanStatus) {
	m.Lock()
	defer m.Unlock()
	m.stats = stats

	// m.dumpStats()

}

func (m *ScannerManager) Close() {
	m.Lock()
	defer m.Unlock()
	m.stats.Running = false

}

func loadLastScanStatus(cached *nutsdb.DB) ScanStatus {

	stats := ScanStatus{}

	cached.View(
		func(tx *nutsdb.Tx) error {

			e, err := tx.Get(`sys`, []byte("last_scan_status"))
			if err != nil {
				return err
			}

			return json.Unmarshal(e.Value, &stats)

		})

	return stats
}
