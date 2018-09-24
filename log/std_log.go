package log

import (
	"fmt"
)

func NewStdLogger(level int) (ILogger, error) {
	return &StdLogger{
		Level: level,
	}, nil
}

type StdLogger struct {
	Level int
}

func (s *StdLogger) SetLevel(level int) {
	s.Level = level
}

func (s *StdLogger) Trace(params ...interface{}) {
	if s.Level <= LogTrace {
		fmt.Printf("%s %c[1;0;33m[TRACE]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *StdLogger) Debug(params ...interface{}) {
	if s.Level <= LogDebug {
		fmt.Printf("%s %c[1;0;34m[DEBUG]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *StdLogger) Info(params ...interface{}) {
	if s.Level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *StdLogger) Error(params ...interface{}) {
	if s.Level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprint(params ...))
	}
}
func (s *StdLogger) Tracef(fromat string, params ...interface{}) {
	if s.Level <= LogTrace {
		fmt.Printf("%s %c[1;0;33m[TRACE]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *StdLogger) Debugf(fromat string, params ...interface{}) {
	if s.Level <= LogDebug {
		fmt.Printf("%s %c[1;0;34m[DEBUG]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *StdLogger) Infof(fromat string, params ...interface{}) {
	if s.Level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}
func (s *StdLogger) Errorf(fromat string, params ...interface{}) {
	if s.Level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR]%c[0m %s\n", Now(), 0x1B, 0x1B, fmt.Sprintf(fromat, params...))
	}
}

func (s *StdLogger) Close() error {
	return nil
}
