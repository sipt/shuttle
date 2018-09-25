package log

import (
	"time"
	"path/filepath"
	"github.com/sipt/shuttle/extension/config"
	"io/ioutil"
)

func init() {
	// path: $HOME/logs
	// level: Debug
	// multiSize: 100MB
	l, err := NewFileLogger(filepath.Join(config.ShuttleHomeDir, "logs"), LogDebug, 100*1000*1000)
	if err != nil {
		ioutil.WriteFile(filepath.Join(config.ShuttleHomeDir, "logs", "error.log"), []byte(err.Error()), 0664)
		panic(err)
	}
	Logger = l
}

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
