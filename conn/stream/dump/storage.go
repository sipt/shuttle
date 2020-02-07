package dump

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/sipt/shuttle/pkg/pool"
	"github.com/sirupsen/logrus"
)

func InitDumpStorage(dir string) error {
	dirPath = path.Join(dir, dirPath)
	return ClearFiles()
}

var (
	dirPath   = "dump"
	allowDump = false
	manager   = &Manager{
		fileMap: make(map[string]io.ReadWriteCloser),
		RWMutex: &sync.RWMutex{},
	}
)

func ReqKey(id int64) string {
	return fmt.Sprintf("req_%d", id)
}

func RespKey(id int64) string {
	return fmt.Sprintf("resp_%d", id)
}

type Manager struct {
	fileMap map[string]io.ReadWriteCloser
	*sync.RWMutex
}

func (m *Manager) InitFiles(id int64) error {
	reqKey := ReqKey(id)
	req, err := CreateFile(dirPath + "/" + reqKey)
	if err != nil {
		return err
	}
	respKey := RespKey(id)
	resp, err := CreateFile(dirPath + "/" + respKey)
	if err != nil {
		req.Close()
		return err
	}
	m.Lock()
	defer m.Unlock()
	m.fileMap[reqKey] = req
	m.fileMap[respKey] = resp
	return nil
}

func InitFiles(id int64) error {
	if !allowDump {
		return nil
	}
	return manager.InitFiles(id)
}

func (m *Manager) GetReqFile(id int64) (io.ReadWriteCloser, bool) {
	m.RLock()
	defer m.RUnlock()
	rw, ok := m.fileMap[ReqKey(id)]
	return rw, ok
}

func (m *Manager) GetRespFile(id int64) (io.ReadWriteCloser, bool) {
	m.RLock()
	defer m.RUnlock()
	rw, ok := m.fileMap[RespKey(id)]
	return rw, ok
}

func (m *Manager) Get(key string) (io.ReadWriteCloser, bool) {
	m.RLock()
	defer m.RUnlock()
	rw, ok := m.fileMap[key]
	return rw, ok
}

func (m *Manager) CloseFiles(id int64) (err error) {
	m.Lock()
	defer m.Unlock()
	key := fmt.Sprintf("req_%d", id)
	rw, ok := m.fileMap[key]
	if ok {
		err = rw.Close()
		delete(m.fileMap, key)
	}
	key = fmt.Sprintf("resp_%d", id)
	rw, ok = m.fileMap[key]
	if ok {
		err = rw.Close()
	}
	return
}

func SaveRequest(id int64, b []byte) (err error) {
	if !allowDump {
		return nil
	}
	data := NewDumpData(ReqKey(id), b)
	buffer <- data
	return
}

func SaveResponse(id int64, b []byte) (err error) {
	if !allowDump {
		return nil
	}
	data := NewDumpData(RespKey(id), b)
	buffer <- data
	return
}

func CloseFiles(id int64) (err error) {
	return manager.CloseFiles(id)
}

func ClearFiles() error {
	isExist, err := PathExists(dirPath)
	if err != nil {
		return err
	}
	if isExist {
		err = os.RemoveAll(dirPath)
		if err != nil {
			return err
		}
	}
	return os.MkdirAll(dirPath, os.ModePerm)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateFile(filePath string) (io.ReadWriteCloser, error) {
	return os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

func NewDumpData(key string, b []byte) *DumpData {
	data := &DumpData{
		key: key,
	}
	data.value = pool.GetBuf()
	if len(b) > len(data.value) {
		data.Release()
	}
	data.value = make([]byte, len(b))
	copy(data.value, b)
	return data
}

type DumpData struct {
	key   string
	value []byte
}

func (d *DumpData) Release() {
	if len(d.value) == 0 && cap(d.value) == 0 {
		return
	}
	pool.PutBuf(d.value)
	d.value = nil
}

var buffer = make(chan *DumpData, 256)

func AutoSave(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-buffer:
			if rw, ok := manager.Get(data.key); ok {
				_, err := rw.Write(data.value)
				if err != nil {
					logrus.WithField("record_id", data.key).WithError(err).Error("[data_dump] save failed")
				}
				data.Release()
			}
		}
	}
}
