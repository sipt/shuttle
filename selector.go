package shuttle

type NewSelector func(group *ServerGroup) (ISelector, error)

var seletors = make(map[string]NewSelector)

func RegisterSelector(method string, newSelector NewSelector) error {
	seletors[method] = newSelector
	return nil
}

func CheckSelector(method string) bool {
	_, ok := seletors[method]
	return ok
}

type ISelector interface {
	Get() (*Server, error)
	Select(string) error
	Refresh() error
	Reset(group *ServerGroup) error
	Destroy()
	Current() IServer
}

func ParseServer(v interface{}) (*ServerGroup, *Server, error) {
	switch v.(type) {
	case *ServerGroup:
		return v.(*ServerGroup), nil, nil
	case *Server:
		return nil, v.(*Server), nil
	}
	return nil, nil, ErrorUnknowType
}
