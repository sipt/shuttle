package close

import "sync"

var closer func(bool) error
var lock = &sync.Mutex{}

type Func func() error

func AppendCloser(cf Func) {
	lock.Lock()
	defer lock.Unlock()
	if closer == nil {
		closer = func(skipErr bool) error {
			return cf()
		}
	}
	closer = func(next func(bool) error) func(bool) error {
		return func(skipErr bool) error {
			err := next(skipErr)
			if !skipErr && err != nil {
				return err
			}
			return cf()
		}
	}(closer)
}

func Close(skipErr bool) error {
	lock.Lock()
	defer lock.Unlock()
	if closer == nil {
		return nil
	}
	err := closer(skipErr)
	closer = nil
	return err
}
