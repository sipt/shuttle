package storage

import (
	"fmt"
	"time"

	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

const DefualtMaxLength = 500

type Record struct {
	ID       int64
	Protocol string
	Created  time.Time
	Proxy    *proxy.Server
	Rule     *rule.Rule
	Status   string
	Up       int
	Down     int
	URL      string
	Dumped   bool
}

type IStorage interface {
	//添加记录到指定的key下
	Put(key string, value Record)
	//当前所有的Key数量
	Keys() []string
	//key下的记录数量
	Count(key string) int
	//获取key的所有记录
	Get(key string) []Record
	// 设置每个Key的上限记录数
	SetLength(l int)
	//清空
	Clear(keys ...string)
}

var storages = make(map[string]IStorage)
var Storage IStorage

func Register(key string, storage IStorage) {
	storages[key] = storage
}

func Use(key string) error {
	var ok bool
	Storage, ok = storages[key]
	if !ok {
		return fmt.Errorf("storage is not support [%s]", key)
	}
	return nil
}
