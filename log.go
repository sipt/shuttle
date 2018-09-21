package shuttle

import (
	"fmt"
	"time"
)

const (
	LogTraceStr = "trace"
	LogDebugStr = "debug"
	LogInfoStr  = "info"
	LogErrorStr = "error"

	LogTrace = 0
	LogDebug = 1
	LogInfo  = 2
	LogError = 3
)

func SetLeve(l string) {
	switch l {
	case LogTraceStr:
		level = LogTrace
	case LogDebugStr:
		level = LogDebug
	case LogInfoStr:
		level = LogInfo
	case LogErrorStr:
		level = LogError
	}
}

var level = LogTrace
var Logger ILogger = &stdLogger{}


func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


