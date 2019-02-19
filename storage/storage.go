package storage

import (
	"fmt"
	"github.com/sipt/shuttle/log"
	"time"

	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

const DefaultStorageCap = 500
const DefaultEngine = EngineMemory

var cap = DefaultStorageCap

type IStorageConfig interface {
	GetStorageEngine() string
	SetStorageEngine(engine string)
	GetStorageCap() int
	SetStorageCap(cap int)
	GetStorageOptions() []string
	SetStorageOptions(options []string)
}

func ApplyConfig(config IStorageConfig) error {
	var e string
	if e = config.GetStorageEngine(); len(e) <= 0 {
		e = DefaultEngine
	}
	cap = config.GetStorageCap()
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
	//获取单条记录
	GetRecord(key string, id int64) *Record
	// 设置每个Key的上限记录数
	SetLength(l int)
	//清空
	Clear(keys ...string)
	//更新Record信息
	Update(key string, id int64, op int, v interface{})
}

type NewStorage func(cap int) IStorage

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
	Storage = f(cap)
	return nil
}

//添加记录到指定的key下
func Put(key string, r Record) {
	if r.Proxy == nil || r.Rule == nil {
		log.Logger.Infof("[Storage] ID:[%d] Policy:[nil] URL:[%s]", r.ID, r.URL)
	} else {
		log.Logger.Infof("[Storage] ID:[%d] Policy:[%s(%s,%s)] URL:[%s]", r.ID, r.Proxy.Name, r.Rule.Type, r.Rule.Value, r.URL)
	}
	Storage.Put(key, r)
}

//当前所有的Key数量
func Keys() []string {
	return Storage.Keys()
}

//key下的记录数量
func Count(key string) int {
	return Storage.Count(key)
}

//获取key的所有记录
func Get(key string) []Record {
	return Storage.Get(key)
}

// 设置每个Key的上限记录数
func SetLength(l int) {
	Storage.SetLength(l)
}

//清空
func Clear(keys ...string) {
	Storage.Clear(keys...)
}

//更新Record信息
func Update(key string, id int64, op int, v interface{}) {
	switch op {
	case RecordUp:
		TrafficUp(v.(int))
	case RecordDown:
		TrafficDown(v.(int))
	}
	Storage.Update(key, id, op, v)
}

//获取单条记录
func GetRecord(key string, id int64) *Record {
	return Storage.GetRecord(key, id)
}
