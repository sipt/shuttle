package shuttle

import (
	"os"
	"fmt"
	"sync"
	"io"
	"github.com/sipt/shuttle/log"
)

var dump IDump

func init() {
	if dump == nil {
		dump = &FileDump{
			Actives: make(map[int64]*SequenceHeap),
		}
	}
	err := dump.Clear()
	if err != nil {
		os.Exit(1)
	}
}

func SetDump(d IDump) {
	dump = d
}

func GetDump() IDump {
	return dump
}

const (
	DumpOrderWrite = iota
	DumpOrderClose

	DumpRequestEntity
	DumpResponseEntity

	DumpFileDir       = "temp"
	DumpRequestFile   = DumpFileDir + string(os.PathSeparator) + "%d_req.dump"
	DumpResponseFile  = DumpFileDir + string(os.PathSeparator) + "%d_resp.dump"
	LargeRequestBody  = 2 * 1024 * 1024 // 5MB
	LargeResponseBody = 2 * 1024 * 1024 // 5MB
)

type IDump interface {
	InitDump(int64) error
	WriteRequest(int64, []byte) (n int, err error)
	WriteResponse(int64, []byte) (n int, err error)
	Dump(int64) (req io.ReadCloser, reqSize int64, resp io.ReadCloser, respSize int64, err error)
	Complete(int64) error
	Clear() error
}

type FileDump struct {
	sync.RWMutex
	Actives      map[int64]*SequenceHeap
	completeList []string
	cancel       chan bool
}

type fileDumpEntity struct {
	data       []byte
	order      int
	entityType int
}

func (f *FileDump) InitDump(id int64) error {
	reqBuf, err := os.OpenFile(fmt.Sprintf(DumpRequestFile, id), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Logger.Errorf("[%d] create data file %s failed: %v", id, err)
		return err
	}
	respBuf, err := os.OpenFile(fmt.Sprintf(DumpResponseFile, id), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Logger.Errorf("[%d] create data file %s failed: %v", id, err)
		return err
	}
	sequenceHeap := NewSequenceHeap()
	f.Lock()
	f.Actives[id] = sequenceHeap
	f.Unlock()
	go func() {
		var data *fileDumpEntity
		for {
			data = sequenceHeap.Pop().(*fileDumpEntity)
			if data == nil {
				reqBuf.Close()
				respBuf.Close()
				return
			}
			switch data.order {
			case DumpOrderWrite:
				switch data.entityType {
				case DumpRequestEntity:
					reqBuf.Write(data.data)
				case DumpResponseEntity:
					respBuf.Write(data.data)
				}
			case DumpOrderClose:
				reqBuf.Close()
				respBuf.Close()
				return
			}
		}
	}()
	return nil
}

func (f *FileDump) WriteRequest(id int64, data []byte) (n int, err error) {
	f.RLock()
	c, ok := f.Actives[id]
	if ok {
		c.Push(&fileDumpEntity{
			data:       data,
			order:      DumpOrderWrite,
			entityType: DumpRequestEntity,
		})
	}
	f.RUnlock()
	return len(data), nil
}
func (f *FileDump) WriteResponse(id int64, data []byte) (n int, err error) {
	f.RLock()
	c, ok := f.Actives[id]
	if ok {
		c.Push(&fileDumpEntity{
			data:       data,
			order:      DumpOrderWrite,
			entityType: DumpResponseEntity,
		})
	}
	f.RUnlock()
	return len(data), nil
}
func (f *FileDump) Dump(id int64) (req io.ReadCloser, reqSize int64, resp io.ReadCloser, respSize int64, err error) {
	file := fmt.Sprintf(DumpRequestFile, id)
	rc, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	req = rc
	info, err := rc.Stat()
	if err != nil {
		return
	}
	reqSize = info.Size()
	file = fmt.Sprintf(DumpResponseFile, id)
	rc, err = os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	resp = rc
	info, err = rc.Stat()
	if err != nil {
		return
	}
	respSize = info.Size()
	return
}
func (f *FileDump) Complete(id int64) error {
	f.RLock()
	_, ok := f.Actives[id]
	f.RUnlock()
	if ok {
		f.Lock()
		c, ok := f.Actives[id]
		if ok {
			delete(f.Actives, id)
		}
		f.Unlock()
		if ok {
			c.Push(&fileDumpEntity{
				order: DumpOrderClose,
			})
			c.Close()
		}
	}
	return nil
}
func (f *FileDump) Clear() error {
	f.Lock()
	for k := range f.Actives {
		c, ok := f.Actives[k]
		if ok {
			c.Push(&fileDumpEntity{
				order: DumpOrderClose,
			})
			c.Close()
		}
	}
	f.Actives = make(map[int64]*SequenceHeap)
	// Clear files
	_, err := os.Stat("temp/")
	if !os.IsNotExist(err) {
		err := os.RemoveAll("temp")
		if err != nil {
			log.Logger.Errorf("delete dir error: %v", err)
			return err
		}
	}
	err = os.Mkdir("temp", os.ModePerm)
	if err != nil {
		log.Logger.Errorf("mkdir failed![%v]\n", err)
		return err
	}
	f.Unlock()
	return nil
}
