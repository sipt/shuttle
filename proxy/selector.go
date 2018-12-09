package proxy

import (
	"errors"
	"fmt"
)

var ErrorUnknowType = errors.New("unknow select type")

type NewSelector func(group *ServerGroup) (ISelector, error)

var seletors = make(map[string]NewSelector)

func RegisterSelector(method string, newSelector NewSelector) {
	seletors[method] = newSelector
}

func GetSelector(method string, group *ServerGroup) (ISelector, error) {
	if s, ok := seletors[method]; ok {
		return s(group)
	}
	return nil, fmt.Errorf("not support select_type [%s]", method)
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
