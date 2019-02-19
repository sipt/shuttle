package storage

import (
	"fmt"
	"time"

	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

const DefaultMaxLength = 500
const DefaultEngine = EngineMemory

type IStorageConfig interface {
	GetStorageEngine() string
	SetStorageEngine(key string)
}

func ApplyConfig(config IStorageConfig) error {
	var e string
	if e = config.GetStorageEngine(); len(e) <= 0 {
		e = DefaultEngine
	}
	return Use(e)
}

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

type NewStorage func() IStorage

var storages = make(map[string]NewStorage)
var Storage IStorage

func Register(key string, f NewStorage) {
	storages[key] = f
}

func Use(key string) error {
	f, ok := storages[key]
	if !ok {
		return fmt.Errorf("storage is not support [%s]", key)
	}
	Storage = f()
	return nil
}
