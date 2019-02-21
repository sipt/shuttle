package storage

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

var Bus = make(chan *Box, 256)

type Box struct {
	//客户端ID，一般是 IP:Port
	ClientID string
	ID       int64
	Op       int
	Value    interface{}
}

func Run() {
	go func() {
		var box *Box
		for {
			box = <-Bus
			switch box.Op {
			case RecordAppend:
				Put(box.ClientID, *(box.Value.(*Record)))
			default:
				Update(box.ClientID, box.ID, box.Op, box.Value)
			}
			go func(box *Box) {
				pusher(box)
			}(box)
		}
	}()
}

type Pusher func(interface{})

var pusher Pusher = func(v interface{}) {} // init empty  pusher
//注册推送
func RegisterPusher(p Pusher) {
	pusher = p
}
