package filter

import (
	"context"
	"time"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/events"

	rulepkg "github.com/sipt/shuttle/rule"
)

var down, up int64 = 0, 0

func init() {
	Register("record", newRecorder)
}
func newRecorder(ctx context.Context, _ map[string]string, filter FilterFunc) (FilterFunc, error) {
	return func(next typ.HandleFunc) typ.HandleFunc {
		handler := filter(next)
		return func(c conn.ICtxConn) {
			req := c.Value(constant.KeyRequestInfo).(typ.RequestInfo)
			rule := c.Value(constant.KeyRule).(*rulepkg.Rule)
			events.Bus <- &events.Event{
				Typ: events.AppendRecordEvent,
				Value: &events.RecordEntity{
					ID:        req.ID(),
					DestAddr:  req.URI(),
					Policy:    rule.String(),
					Status:    events.ActiveStatus,
					Timestamp: time.Now(),
					Protocol:  c.Value(constant.KeyProtocol).(string),
				},
			}
			handler(c)
		}
	}, nil
}
