package shuttle

import (
	"time"
	"sync"
)

const (
	RecordStatus = 1
	RecordUp     = 2
	RecordDown   = 3

	RecordStatusActive    = "Active"
	RecordStatusCompleted = "Completed"
)

var maxCount = 1000
var recordChan chan *Record
var storage *LinkedList

func init() {
	recordChan = make(chan *Record, 16)
	storage = &LinkedList{}
	go func() {
		for {
			storage.Append(<-recordChan)
		}
	}()
}

func GetRecords() []Record {
	return storage.List()
}
func ClearRecords() {
	storage.Clear()
}
func GetRecord(id int64) *Record {
	return storage.Get(id)
}

type Record struct {
	ID       int64
	Protocol string
	Created  time.Time
	Proxy    *Server
	Rule     *Rule
	Status   string
	Up       int
	Down     int
	URL      string
}

type LinkedList struct {
	sync.RWMutex
	head, tail *node
	count      int
}
type node struct {
	record *Record
	next   *node
	sync.RWMutex
}

func (l *LinkedList) Get(id int64) *Record {
	l.RLock()
	index := l.head
	for index != nil {
		if index.record.ID == id {
			l.RUnlock()
			return index.record
		}
		index = index.next
	}
	l.RUnlock()
	return nil
}
func (l *LinkedList) List() []Record {
	if l.count == 0 {
		return []Record{}
	}
	l.Lock()
	list := make([]Record, l.count)
	index := l.head
	for i := range list {
		list[i] = *index.record
		index = index.next
	}
	l.Unlock()
	return list
}
func (l *LinkedList) Append(r *Record) {
	l.Lock()
	Logger.Debugf("[Storage] Policy:[%s(%s,%s)] URL:[%s]", r.Proxy.Name, r.Rule.Type, r.Rule.Value, r.URL)
	if l.head == nil {
		l.head = &node{record: r}
		l.tail = l.head
	} else {
		l.tail.next = &node{record: r}
		l.tail = l.tail.next
	}
	l.count ++

	for l.count > maxCount {
		// 收缩
		l.head.next, l.head = nil, l.head.next
		l.count --
	}
	l.Unlock()
}
func (l *LinkedList) Put(id int64, op int, v interface{}) {
	index := l.head
	for index != nil {
		if index.record.ID == id {
			index.Put(op, v)
		}
		index = index.next
	}
}
func (l *LinkedList) Clear() {
	l.Lock()
	l.head = nil
	l.count = 0
	l.Unlock()
}
func (n *node) Put(op int, v interface{}) {
	n.Lock()
	switch op {
	case RecordStatus:
		n.record.Status = v.(string)
	case RecordUp:
		n.record.Up += v.(int)
	case RecordDown:
		n.record.Down += v.(int)
	}
	n.Unlock()
}
