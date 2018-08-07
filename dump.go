package shuttle

import (
	"io"
	"os"
	"fmt"
	"io/ioutil"
	"sync"
)

var dump IDump

func init() {
	if dump == nil {
		dump = &FileDump{
			Actives: make(map[string]io.WriteCloser),
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

type IDump interface {
	WriteRequest(int64, []byte) (n int, err error)
	WriteResponse(int64, []byte) (n int, err error)
	ReadRequest(int64) ([]byte, error)
	ReadResponse(int64) ([]byte, error)
	Complete(int64) error
	Clear() error
}

type FileDump struct {
	sync.RWMutex
	Actives      map[string]io.WriteCloser
	OldActives   map[string]io.WriteCloser
	completeList []string
}

func (f *FileDump) WriteRequest(id int64, data []byte) (n int, err error) {
	name := fmt.Sprintf("./temp/%d_request.txt", id)
	return f.write(name, data)
}
func (f *FileDump) WriteResponse(id int64, data []byte) (n int, err error) {
	name := fmt.Sprintf("./temp/%d_reponse.txt", id)
	return f.write(name, data)
}
func (f *FileDump) write(fileName string, data []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()
	if len(f.OldActives) > 0 {
		_, ok := f.OldActives[fileName]
		if ok {
			return
		}
	}
	w, ok := f.Actives[fileName]
	if !ok {
		Logger.Debugf("[DUMP FILE] create file: %s", fileName)
		file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return 0, err
		}
		f.Actives[fileName] = file
		w = file
	}
	Logger.Debugf("[DUMP FILE] write to file: %s", fileName)
	return w.Write(data)
}
func (f *FileDump) ReadRequest(id int64) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("./temp/%d_request.txt", id))
}
func (f *FileDump) ReadResponse(id int64) ([]byte, error) {
	return ioutil.ReadFile(fmt.Sprintf("./temp/%d_reponse.txt", id))
}
func (f *FileDump) Complete(id int64) error {
	f.close(fmt.Sprintf("./temp/%d_request.txt", id))
	f.close(fmt.Sprintf("./temp/%d_reponse.txt", id))
	return nil
}
func (f *FileDump) close(name string) error {
	f.Lock()
	defer f.Unlock()
	if len(f.OldActives) > 0 {
		w, ok := f.OldActives[name]
		if ok {
			delete(f.OldActives, name)
			w.Close()
			return nil
		}
	}
	w, ok := f.Actives[name]
	if ok {
		delete(f.Actives, name)
		w.Close()
	}
	return nil
}
func (f *FileDump) Clear() error {
	f.Lock()
	defer f.Unlock()
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
	if len(f.Actives) > 0 {
		f.OldActives = f.Actives
		f.Actives = make(map[string]io.WriteCloser)
	}
	return nil
}
