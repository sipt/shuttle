package log

import "time"

var Logger ILogger = &StdLogger{Level: LogDebug}

func SetLogger(logger ILogger) {
	Logger = logger
}

const (
	LogTrace = 0
	LogDebug = 1
	LogInfo  = 2
	LogError = 3
)

var LevelMap = map[string]int{
	"trace": LogTrace,
	"debug": LogDebug,
	"info":  LogInfo,
	"error": LogError,
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

type ILogger interface {
	SetLevel(int)
	Trace(...interface{})
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
	Tracef(string, ...interface{})
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Close() error
}
