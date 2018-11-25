package log

import (
	"errors"
	"time"
)

const (
	LogModeOff     = "off"
	LogModeConsole = "console"
	LogModeFile    = "file"
)

type ILogConfig interface {
	GetLogLevel() string
}

func InitLogger(logMode, logPath string, logConfig ILogConfig) (err error) {
	var (
		ok        bool
		levelFlag int
		level     = logConfig.GetLogLevel()
	)
	switch logMode {
	case LogModeOff:
		Logger, err = NewSkipLogger()
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	case LogModeConsole:
		levelFlag, ok = LevelMap[level]
		if !ok {
			return errors.New("not support LogLevel:" + level)
		}
		Logger, err = NewStdLogger(levelFlag)
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	case LogModeFile:
		levelFlag, ok = LevelMap[level]
		if !ok {
			return errors.New("not support LogLevel:" + level)
		}
		//multiSize: 100MB
		Logger, err = NewFileLogger(logPath, levelFlag, 100*1000*1000)
		if err != nil {
			return errors.New("init logger failed:" + err.Error())
		}
	default:
		return errors.New("not support LogMode:" + logMode)
	}
	return
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
