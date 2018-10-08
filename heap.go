package shuttle

import (
	"container/heap"
	"sync"
	"sync/atomic"
)

type Item struct {
	Value  interface{}
	Ticket int64
}

type MinHeap []*Item

func (m *MinHeap) Len() int { return len(*m) }

func (m *MinHeap) Less(i, j int) bool {
	return (*m)[i].Ticket < (*m)[j].Ticket
}

func (m *MinHeap) Swap(i, j int) {
	(*m)[i], (*m)[j] = (*m)[j], (*m)[i]
}

func (m *MinHeap) Push(x interface{}) {
	item := x.(*Item)
	*m = append(*m, item)
}

func (m *MinHeap) Pop() interface{} {
	old := *m
	n := len(old)
	item := old[n-1]
	*m = old[:n-1]
	return item
}

func NewMinArrange() *MinArrange {
	h := make(MinHeap, 0, 8)
	return &MinArrange{
		oldTicket: 0,
		minHeap:   &h,
		result:    make(chan *Item, 32),
		closed:    false,
		closeChan: make(chan bool, 1),
	}
}

type MinArrange struct {
	oldTicket int64
	minHeap   *MinHeap
	result    chan *Item
	closeChan chan bool
	closed    bool
	sync.Mutex
}

func (m *MinArrange) Pop() (item *Item) {
	item = <-m.result
	return
}

func (m *MinArrange) Push(v *Item) {
	m.Lock()
	heap.Push(m.minHeap, v)
Loop:
	for len(*m.minHeap) > 0 && m.oldTicket+1 == (*m.minHeap)[0].Ticket {
		select {
		case m.result <- (*m.minHeap)[0]:
			heap.Pop(m.minHeap)
			m.oldTicket ++
		default:
			break Loop
		}
	}
	m.Unlock()
}

func (m *MinArrange) Close() {
	m.closed = true
	m.closeChan <- true
}

func NewSequenceHeap() *SequenceHeap {
	return &SequenceHeap{
		NewMinArrange(), 0,
	}
}

type SequenceHeap struct {
	*MinArrange
	ticket int64
}

func (s *SequenceHeap) Pop() interface{} {
	v := s.MinArrange.Pop()
	if v != nil {
		return v.Value
	}
	return nil
}

func (s *SequenceHeap) Push(v interface{}) {
	s.MinArrange.Push(&Item{
		Value:  v,
		Ticket: s.getTicket(),
	})
}

func (s *SequenceHeap) getTicket() int64 {
	return atomic.AddInt64(&s.ticket, 1)
}
