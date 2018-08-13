package shuttle

import (
	"os"
	"fmt"
	"io/ioutil"
	"sync"
)

var dump IDump

func init() {
	if dump == nil {
		dump = &FileDump{
			Actives: make(map[int64]chan *fileDumpEntity),
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
)

type IDump interface {
	InitDump(int64) error
	WriteRequest(int64, []byte) (n int, err error)
	WriteResponse(int64, []byte) (n int, err error)
	ReadRequest(int64) ([]byte, error)
	ReadResponse(int64) ([]byte, error)
	Complete(int64) error
	Clear() error
}

type FileDump struct {
	sync.RWMutex
	Actives      map[int64]chan *fileDumpEntity
	completeList []string
	cancel       chan bool
}

type fileDumpEntity struct {
	data       []byte
	order      int
	entityType int
}

func (f *FileDump) InitDump(id int64) error {
	request, err := os.OpenFile(fmt.Sprintf("./temp/%d_request.txt", id), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	response, err := os.OpenFile(fmt.Sprintf("./temp/%d_reponse.txt", id), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	dataChan := make(chan *fileDumpEntity, 8)
	f.Actives[id] = dataChan
	go func() {
		var data *fileDumpEntity
		for {
			data = <-dataChan
			switch data.order {
			case DumpOrderWrite:
				switch data.entityType {
				case DumpRequestEntity:
					request.Write(data.data)
				case DumpResponseEntity:
					response.Write(data.data)
				}
			case DumpOrderClose:
				request.Close()
				response.Close()
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
		c <- &fileDumpEntity{
			data:       data,
			order:      DumpOrderWrite,
			entityType: DumpRequestEntity,
		}
	}
	f.RUnlock()
	return len(data), nil
}
func (f *FileDump) WriteResponse(id int64, data []byte) (n int, err error) {
	f.RLock()
	c, ok := f.Actives[id]
	if ok {
		c <- &fileDumpEntity{
			data:       data,
			order:      DumpOrderWrite,
			entityType: DumpResponseEntity,
		}
	}
	f.RUnlock()
	return len(data), nil
}
func (f *FileDump) ReadRequest(id int64) ([]byte, error) {
	file := fmt.Sprintf("./temp/%d_request.txt", id)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	return ioutil.ReadFile(fmt.Sprintf(file, id))
}
func (f *FileDump) ReadResponse(id int64) ([]byte, error) {
	file := fmt.Sprintf("./temp/%d_reponse.txt", id)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	return ioutil.ReadFile(fmt.Sprintf(file, id))
}
func (f *FileDump) Complete(id int64) error {
	_, ok := f.Actives[id]
	if ok {
		f.Lock()
		c, ok := f.Actives[id]
		if ok {
			c <- &fileDumpEntity{
				order: DumpOrderClose,
			}
			delete(f.Actives, id)
		}
		f.Unlock()
	}
	return nil
}
func (f *FileDump) Clear() error {
	f.Lock()
	for k := range f.Actives {
		c, ok := f.Actives[k]
		if ok {
			c <- &fileDumpEntity{
				order: DumpOrderClose,
			}
		}
	}
	f.Actives = make(map[int64]chan *fileDumpEntity)
	// Clear files
	_, err := os.Stat("temp/")
	if !os.IsNotExist(err) {
		err := os.RemoveAll("temp")
		if err != nil {
			Logger.Errorf("delete dir error: %v", err)
			return err
		}
	}
	err = os.Mkdir("temp", os.ModePerm)
	if err != nil {
		Logger.Errorf("mkdir failed![%v]\n", err)
		return err
	}
	f.Unlock()
	return nil
}
