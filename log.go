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

type ILogger interface {
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
	Tracef(string, ...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type stdLogger struct{}

func (s *stdLogger) Trace(params ...interface{}) {
	if level <= LogTrace {
		fmt.Printf("%s %c[1;0;33m[TRACE]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *stdLogger) Debug(params ...interface{}) {
	if level <= LogDebug {
		fmt.Printf("%s %c[1;0;34m[DEBUG]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *stdLogger) Info(params ...interface{}) {
	if level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *stdLogger) Error(params ...interface{}) {
	if level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *stdLogger) Tracef(fromat string, params ...interface{}) {
	if level <= LogTrace {
		fmt.Printf("%s %c[1;0;33m[TRACE]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *stdLogger) Debugf(fromat string, params ...interface{}) {
	if level <= LogDebug {
		fmt.Printf("%s %c[1;0;34m[DEBUG]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *stdLogger) Infof(fromat string, params ...interface{}) {
	if level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *stdLogger) Errorf(fromat string, params ...interface{}) {
	if level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
