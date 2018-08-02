package shuttle

import "errors"

var (
	ErrorReadTimeOut  = errors.New("read time out")
	ErrorWriteTimeOut = errors.New("write time out")
	ErrorReject       = errors.New("connection reject")

	ErrorServerNotFound = errors.New("server or server group not found")

	ErrorUnknowType = errors.New("unknow type")
)
