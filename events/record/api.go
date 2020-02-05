package record

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/events"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/records", recordsHandleFunc)
	e.DELETE("/api/records", clearRecordsHandleFunc)

	e.GET("/ws/records/events", recordEventsHandleFunc)
}

func recordsHandleFunc(c *gin.Context) {
	list := make([]RecordEntity, 0, recordStarge.Len())
	recordStarge.Range(func(v interface{}) bool {
		r := v.(*RecordEntity)
		list = append(list, *r)
		return false
	})
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: list,
	})
}

func clearRecordsHandleFunc(c *gin.Context) {
	recordStarge.Clear()
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
	})
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	wsConnID  int64 = 0
	wsConnMap       = make(map[int64]*websocket.Conn)
)

func recordEventsHandleFunc(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	id := atomic.AddInt64(&wsConnID, 1)
	go func() {
		defer conn.Close()
		for {
			typ, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if typ == websocket.CloseMessage {
				delete(wsConnMap, id)
				return
			}
		}
	}()
	wsConnMap[id] = conn
}

func notifyClient(eventType events.EventType, r *RecordEntity) {
	if len(wsConnMap) == 0 {
		return
	}
	resp := &EventResponse{
		Typ: eventType,
		Value: &RecordResponse{
			ID:       r.ID,
			DestAddr: r.DestAddr,
			Policy:   r.Policy,
			Up:       r.Up,
			Down:     r.Down,
			Status:   r.Status,
			Protocol: r.Protocol,
			Duration: r.Duration.Nanoseconds() / int64(time.Millisecond),
			Conn:     r.Conn,
		},
	}
	if !r.Timestamp.IsZero() {
		resp.Value.Timestamp = r.Timestamp.UnixNano() / int64(time.Millisecond)
	}
	for _, v := range wsConnMap {
		_ = v.WriteJSON(resp)
	}
}

type EventResponse struct {
	Typ   events.EventType `json:"typ"`
	Value *RecordResponse  `json:"value"`
}
type RecordResponse struct {
	ID        int64        `json:"id,omitempty"`
	DestAddr  string       `json:"dest_addr,omitempty"`
	Policy    string       `json:"policy,omitempty"`
	Up        int64        `json:"up"`
	Down      int64        `json:"down"`
	Status    RecordStatus `json:"status,omitempty"`
	Timestamp int64        `json:"timestamp,omitempty"`
	Protocol  string       `json:"protocol,omitempty"`
	Duration  int64        `json:"duration,omitempty"`
	Conn      *ConnEntity  `json:"conn,omitempty"`
}
