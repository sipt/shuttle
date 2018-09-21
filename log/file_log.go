package log

import (
	"fmt"
	"os"
	"sync"
	"path/filepath"
	"time"
	"io"
)

func NewFileLogger(filePath string, level, multiSize int) (ILogger, error) {
	l := &Logger{
		Level: level,
	}
	lf := &LogFile{
		path:      filePath,
		multiSize: multiSize,
	}
	return l, nil
}

type Logger struct {
	Out   io.WriteCloser
	Level int
}

func (l *Logger) Trace(params ...interface{}) {
	if l.Level <= LogTrace {
		fmt.Printf("%s [TRACE] %s\n", Now(), fmt.Sprint(params ...))
	}
}
func (l *Logger) Debug(params ...interface{}) {
	if l.Level <= LogDebug {
		fmt.Printf("%s [DEBUG] %s\n", Now(), fmt.Sprint(params ...))
	}
}
func (l *Logger) Info(params ...interface{}) {
	if l.Level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO] %s\n", Now(), fmt.Sprint(params ...))
	}
}
func (l *Logger) Error(params ...interface{}) {
	if l.Level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR] %s\n", Now(), fmt.Sprint(params ...))
	}
}
func (l *Logger) Tracef(format string, params ...interface{}) {
	if l.Level <= LogTrace {
		fmt.Printf("%s %c[1;0;33m[TRACE] %s\n", Now(), fmt.Sprintf(format, params...))
	}
}
func (l *Logger) Debugf(format string, params ...interface{}) {
	if l.Level <= LogDebug {
		fmt.Printf("%s [DEBUG] %s\n", Now(), fmt.Sprintf(format, params...))
	}
}
func (l *Logger) Infof(format string, params ...interface{}) {
	if l.Level <= LogInfo {
		fmt.Printf("%s %c[1;0;32m[INFO] %s\n", Now(), fmt.Sprintf(format, params...))
	}
}
func (l *Logger) Errorf(format string, params ...interface{}) {
	if l.Level <= LogError {
		fmt.Printf("%s %c[1;0;31m[ERROR] %s\n", Now(), fmt.Sprintf(format, params...))
	}
}

type LogFile struct {
	path      string
	file      *os.File
	multiSize int
	size      int
	sync.RWMutex
}

func (l *LogFile) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	n, err := l.file.Write(b)
	l.size += n
	if l.size > l.multiSize {
		l.file.Close()
		l.Create()
		l.size = 0
	}
	return n, err
}

func (l *LogFile) Create() error {
	file, err := os.OpenFile(filepath.Join(l.path, time.Now().Format("2006-01-02_15:04:05")+".log"),
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}

func (l *LogFile) Close() error {
	l.Lock()
	defer l.Unlock()
	return l.file.Close()
}
