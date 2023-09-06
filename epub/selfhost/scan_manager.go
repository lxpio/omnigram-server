package selfhost

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/nexptr/omnigram-server/conf"
	"github.com/nexptr/omnigram-server/epub/schema"
	"github.com/nexptr/omnigram-server/log"
	"github.com/nexptr/omnigram-server/utils"
)

// const statsCachePath = `config`
const statsfile = `state.json`

// ScanStatus 扫描状态
type ScanStatus struct {
	Running bool     `json:"running"`
	Count   int      `json:"count"`
	Errs    []string `json:"errs"`
}

type ScannerManager struct {
	cf *conf.Config

	kv    schema.KV
	store *schema.Store

	ctx context.Context

	sync.RWMutex
	scan *Scanner

	stats ScanStatus
}

func NewScannerManager(ctx context.Context, cf *conf.Config, kv schema.KV, store *schema.Store) (*ScannerManager, error) {

	scanner := &ScannerManager{
		cf:    cf,
		kv:    kv,
		store: store,
		ctx:   ctx,
	}

	//获取本地存储的状态
	if obj, err := kv.GetObject(ctx, utils.ConfigBucket, statsfile); err == nil {
		json.Unmarshal(obj.Data, &scanner.stats)
	}

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

func (m *ScannerManager) Start(maxThread int) {

	if m.IsRunning() {
		log.E(`扫描器已经在执行，放弃执行`)
		return
	}
	log.I(`启动文件目录扫描`)
	m.newScan(m.cf.EpubOptions.DataPath, maxThread)

}

func (m *ScannerManager) newScan(path string, maxThread int) {
	m.Lock()
	m.scan = NewScan(path, maxThread) //new scanner
	m.stats.Running = true
	m.Unlock()

	m.scan.Start(m, maxThread)
}

func (m *ScannerManager) Stop() {

	m.Lock()
	defer m.Unlock()

	if m.scan != nil {
		m.scan.Stop()
	}

}

func (m *ScannerManager) dumpStats() error {

	bytes, _ := json.Marshal(m.stats)

	obj := &schema.Object{
		Key:          statsfile,
		Size:         0,
		LastModified: time.Time{},
		Data:         bytes,
	}

	return m.kv.PutObject(m.ctx, utils.ConfigBucket, obj)
}

func (m *ScannerManager) updateStatus(stats ScanStatus) {
	m.Lock()
	defer m.Unlock()
	m.stats = stats

	m.dumpStats()

}

func (m *ScannerManager) Close() {
	m.Stop()
}
