package constant

var (
	EventShutdown          = &EventObj{Type: 0}
	EventReloadConfig      = &EventObj{Type: 1}
	EventRestartHttpProxy  = &EventObj{Type: 2}
	EventRestartSocksProxy = &EventObj{Type: 3}
	EventRestartController = &EventObj{Type: 4}
	EventUpgrade           = &EventObj{Type: 5}
)

type EventObj struct {
	Type int
	data interface{}
}

func (e *EventObj) SetData(data interface{}) *EventObj {
	eb := &EventObj{}
	*eb = *e
	eb.data = data
	return eb
}

func (e *EventObj) GetData() interface{} {
	return e.data
}
