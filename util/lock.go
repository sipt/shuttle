package util

import "sync"

var lockMap = &sync.Map{}

func getLock(key interface{}) *sync.RWMutex {
	v, ok := lockMap.Load(key)
	if !ok {
		v = &sync.RWMutex{}
		lockMap.Store(key, v)
	}
	return v.(*sync.RWMutex)
}

func RLock(key string) {
	lock := getLock(key)
	lock.RLock()
}

func RUnLock(key string) {
	v, ok := lockMap.Load(key)
	if ok {
		lock := v.(*sync.RWMutex)
		lock.RUnlock()
	}
}
func Lock(key string) {
	lock := getLock(key)
	lock.Lock()
}

func UnLock(key string) {
	v, ok := lockMap.Load(key)
	if ok {
		lock := v.(*sync.RWMutex)
		lock.Unlock()
	}
}
