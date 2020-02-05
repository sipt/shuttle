package record

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/events"
	"github.com/sipt/shuttle/pkg/storage"
	"github.com/sirupsen/logrus"
)

const (
	AppendRecordEvent       events.EventType = 1
	UpdateRecordUpEvent     events.EventType = 2
	UpdateRecordDownEvent   events.EventType = 3
	UpdateRecordStatusEvent events.EventType = 4
)

func init() {
	// append record
	events.RegisterEvent(AppendRecordEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		AppendRecord(ctx, r)
		notifyClient(AppendRecordEvent, r)
		return nil
	})
	// update record up
	events.RegisterEvent(UpdateRecordUpEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			atomic.AddInt64(&re.Up, r.Up)
			r.Up = re.Up
		})
		notifyClient(UpdateRecordUpEvent, r)
		return nil
	})
	// update record down
	events.RegisterEvent(UpdateRecordDownEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			atomic.AddInt64(&re.Down, r.Down)
			r.Down = re.Down
		})
		notifyClient(UpdateRecordDownEvent, r)
		return nil
	})
	// update record status
	events.RegisterEvent(UpdateRecordStatusEvent, func(ctx context.Context, v interface{}) error {
		r, ok := v.(*RecordEntity)
		if !ok {
			return errors.Errorf("[%s] is not RecordEntity", reflect.TypeOf(v).Kind().String())
		}
		UpdateRecord(ctx, r.ID, func(re *RecordEntity) {
			re.Status = r.Status
		})
		notifyClient(UpdateRecordDownEvent, r)
		return nil
	})
}

type ConnEntity struct {
	ID         int64  `json:"id"`
	SourceAddr string `json:"source_addr"`
	DestAddr   string `json:"dest_addr"`
}

type RecordStatus string

func (r *RecordStatus) String() string {
	return string(*r)
}

var (
	ActiveStatus    RecordStatus = "Active"
	CompletedStatus RecordStatus = "Completed"
	FailedStatus    RecordStatus = "Failed"
	RejectedStatus  RecordStatus = "Rejected"
)

type RecordEntity struct {
	ID        int64         `json:"id,omitempty"`
	DestAddr  string        `json:"dest_addr,omitempty"`
	Policy    string        `json:"policy,omitempty"`
	Up        int64         `json:"up"`
	Down      int64         `json:"down"`
	Status    RecordStatus  `json:"status,omitempty"`
	Timestamp time.Time     `json:"timestamp,omitempty"`
	Protocol  string        `json:"protocol,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
	Conn      *ConnEntity   `json:"conn,omitempty"`
}

var recordStarge = storage.NewLRUList(1000)

func AppendRecord(ctx context.Context, record *RecordEntity) {
	logrus.WithField("recode", *record).Debug("events: append record")
	recordStarge.PushBack(record)
}
func RangeRecord(ctx context.Context, f func(*RecordEntity) bool) {
	recordStarge.Range(func(v interface{}) bool {
		return f(v.(*RecordEntity))
	})
}
func UpdateRecord(ctx context.Context, id int64, f func(*RecordEntity)) {
	recordStarge.Range(func(v interface{}) bool {
		if r, ok := v.(*RecordEntity); ok && r.ID == id {
			f(r)
			return true
		}
		return false
	})
}
