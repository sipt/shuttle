package shuttle

import (
	"os"
	"fmt"
	"io/ioutil"
	"sync"
	"bytes"
	"github.com/sipt/shuttle/pool"
	"bufio"
	"net/http"
	"compress/gzip"
	"compress/zlib"
	"io"
	"encoding/base64"
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
	Dump(int64) ([]byte, error)
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
	reqBuf := bytes.NewBuffer(pool.GetBuf()[:0])
	respBuf := bytes.NewBuffer(pool.GetBuf()[:0])
	dataChan := make(chan *fileDumpEntity, 8)
	f.Lock()
	f.Actives[id] = dataChan
	f.Unlock()
	go func() {
		var data *fileDumpEntity
		for {
			data = <-dataChan
			switch data.order {
			case DumpOrderWrite:
				switch data.entityType {
				case DumpRequestEntity:
					reqBuf.Write(data.data)
				case DumpResponseEntity:
					respBuf.Write(data.data)
				}
			case DumpOrderClose:
				file, err := os.OpenFile(fmt.Sprintf("./temp/%d_data.txt", id), os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					Logger.Errorf("[%d] create data file failed: %v", id, err)
				}
				file.WriteString(`{"req":"`)
				file.WriteString(base64.StdEncoding.EncodeToString(reqBuf.Bytes()))
				file.WriteString(`"`)
				//解析返回值
				b := bufio.NewReader(reqBuf)
				req, err := http.ReadRequest(b)
				if err != nil {
					Logger.Errorf("[%d] parse http request failed: %v", id, err)
				}
				b = bufio.NewReader(respBuf)
				resp, err := http.ReadResponse(b, req)
				if err != nil {
					Logger.Errorf("[%d] parse http response failed: %v", id, err)
				}
				var r io.Reader
				if resp.Header.Get("Content-Encoding") == "gzip" {
					r, err = gzip.NewReader(resp.Body)
					if err != nil {
						Logger.Errorf("[%d] gzip init for response failed: %v", id, err)
					}
				} else if resp.Header.Get("Content-Encoding") == "deflate" {
					r, err = zlib.NewReader(resp.Body)
					if err != nil {
						Logger.Errorf("[%d] deflate init for response failed: %v", id, err)
					}
				} else {
					r = resp.Body
				}

				file.WriteString(`,"resp_body":"`)
				reqBuf.Reset()
				reqBuf.ReadFrom(r)
				file.WriteString(base64.StdEncoding.EncodeToString(reqBuf.Bytes()))
				resp.Body.Close()
				file.WriteString(`","resp_header":"`)
				reqBuf.Reset()
				resp.Write(reqBuf)
				file.WriteString(base64.StdEncoding.EncodeToString(reqBuf.Bytes()))
				file.WriteString(`"}`)
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
func (f *FileDump) Dump(id int64) ([]byte, error) {
	file := fmt.Sprintf("./temp/%d_data.txt", id)
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return []byte{}, nil
	}
	return ioutil.ReadFile(file)
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
			c <- &fileDumpEntity{
				order: DumpOrderClose,
			}
		}
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
