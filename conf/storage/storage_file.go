package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	KeyFile = "file"
)

func init() {
	RegisterStorage(KeyFile, newFileStorage)
}

func newFileStorage(params map[string]string) (IStorage, error) {
	path, ok := params["path"]
	if !ok {
		return nil, errors.Errorf("not found [path] in params")
	}
	if len(path) == 0 {
		return nil, errors.Errorf("[path] is empty in params")
	}

	return &fileStorage{
		filePath: path,
		RWMutex:  new(sync.RWMutex),
	}, nil
}

type fileStorage struct {
	filePath string
	notify   func()
	*sync.RWMutex
}

// Load: load config from disk? As JSON? YAML? TOML?
func (f *fileStorage) Load() ([]byte, error) {
	f.RLock()
	defer f.RUnlock()
	data, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "[Load] file failed: %s", f.filePath)
	}
	return data, nil
}

// Save: save config to file? upload to server?
func (f *fileStorage) Save(data []byte) error {
	f.Lock()
	defer f.Unlock()
	err := ioutil.WriteFile(f.filePath, data, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "[Save] file failed: %s", f.filePath)
	}
	return nil
}

// RegisterNotify
func (f *fileStorage) RegisterNotify(ctx context.Context, notify func()) error {
	if notify == nil {
		return nil
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrapf(err, "[RegisterNotify] failed: %s", f.filePath)
	}

	go func() {
		log := logrus.WithField("method", "FileWatcher")
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				log.Infof("modified file: %s", event.String(), ok)
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					notify()
					log.Infof("modified file: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("error: %s", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	err = watcher.Add(f.filePath)
	if err != nil {
		return errors.Wrapf(err, "[RegisterNotify] watcher file failed: %s", f.filePath)
	}
	fmt.Println("register notify success: ", f.filePath)
	return nil
}
