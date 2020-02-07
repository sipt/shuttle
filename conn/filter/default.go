package filter

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sipt/shuttle/conn/stream/dump"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/events"
	"github.com/sipt/shuttle/events/record"

	rulepkg "github.com/sipt/shuttle/rule"
)

var down, up int64 = 0, 0

func init() {
	Register("record", newRecorder)
}
func newRecorder(ctx context.Context, _ map[string]string, next typ.HandleFunc) (typ.HandleFunc, error) {
	return func(c conn.ICtxConn) {
		req := c.Value(constant.KeyRequestInfo).(typ.RequestInfo)
		rule := c.Value(constant.KeyRule).(*rulepkg.Rule)
		events.Bus <- &events.Event{
			Typ: record.AppendRecordEvent,
			Value: &record.RecordEntity{
				ID:        req.ID(),
				DestAddr:  req.URI(),
				Policy:    rule.String(),
				Status:    record.ActiveStatus,
				Timestamp: time.Now(),
				Protocol:  c.Value(constant.KeyProtocol).(string),
			},
		}
		err := dump.InitFiles(c, req.ID())
		if err != nil {
			logrus.WithField("record_id", req.ID()).WithError(err).Error("[data_dump] init files failed")
		}
		next(c)
	}, nil
}
