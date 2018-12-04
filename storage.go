package shuttle

import (
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
	"sync"
	"sync/atomic"
	"time"
)

const (
	RecordStatus = 1
	RecordUp     = 2
	RecordDown   = 3
	RecordAppend = 4
	RecordRemove = 5

	RecordStatusActive    = "Active"
	RecordStatusCompleted = "Completed"
	RecordStatusReject    = "Reject"
	RecordStatusFailed    = "Failed"
)

type Pusher func(interface{})

var maxCount = 500
var boxChan chan *Box
var storage *LinkedList
var speed *Speed
var pusher Pusher = func(v interface{}) {} // init empty  pusher

func init() {
	boxChan = make(chan *Box, 64)
	storage = &LinkedList{}
	speed = &Speed{
		Cancel: make(chan bool, 1),
	}
	speed.Start()
	go func() {
		var box *Box
		for {
			box = <-boxChan
			switch box.Op {
			case RecordAppend:
				storage.Append(box.Value.(*Record))
			default:
				storage.Put(box.ID, box.Op, box.Value)
			}
			go func(box *Box) {
				pusher(box)
			}(box)
		}
	}()

	//init traffic
	conn.InitTrafficChannel(func(recordID int64, n int) {
		boxChan <- &Box{recordID, RecordUp, n}
	}, func(recordID int64, n int) {
		boxChan <- &Box{recordID, RecordDown, n}
	})
}

//注册推送
func RegisterPusher(p Pusher) {
	pusher = p
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
func CurrentSpeed() (int, int) {
	return speed.UpSpeed, speed.DownSpeed
}

type Speed struct {
	UpSpeed   int
	DownSpeed int
	UpBytes   int
	DownBytes int
	Cancel    chan bool
	status    int32
}

func (s *Speed) Start() {
	if atomic.CompareAndSwapInt32(&s.status, 0, 1) {
		go func() {
			t := time.NewTicker(time.Second)
			for {
				select {
				case <-t.C:
					s.UpSpeed, s.UpBytes = s.UpBytes, 0
					s.DownSpeed, s.DownBytes = s.DownBytes, 0
				case <-s.Cancel:
					return
				}
			}
		}()
	}
}

type Box struct {
	ID    int64
	Op    int
	Value interface{}
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
	if r.Proxy == nil || r.Rule == nil {
		log.Logger.Infof("[Storage] ID:[%d] Policy:[nil] URL:[%s]", r.ID, r.URL)
	} else {
		log.Logger.Infof("[Storage] ID:[%d] Policy:[%s(%s,%s)] URL:[%s]", r.ID, r.Proxy.Name, r.Rule.Type, r.Rule.Value, r.URL)
	}
	if l.head == nil {
		l.head = &node{record: r}
		l.tail = l.head
	} else {
		l.tail.next = &node{record: r}
		l.tail = l.tail.next
	}
	l.count ++

	for l.count > maxCount {
		go func(id int64) {
			pusher(&Box{
				Op:    RecordRemove,
				Value: id,
			})
		}(l.head.record.ID)
		// 收缩
		l.head.next, l.head = nil, l.head.next
		l.count --
	}
	l.Unlock()
}
func (l *LinkedList) Put(id int64, op int, v interface{}) {
	l.RLock()
	index := l.head
	for index != nil {
		if index.record.ID == id {
			index.Put(op, v)
			break
		}
		index = index.next
	}
	l.RUnlock()
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
		s := v.(int)
		n.record.Up += s
		if speed != nil {
			speed.UpBytes += s
		}
	case RecordDown:
		s := v.(int)
		n.record.Down += s
		if speed != nil {
			speed.DownBytes += s
		}
	}
	n.Unlock()
}
