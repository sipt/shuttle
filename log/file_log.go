package log

import (
	"fmt"
	"os"
	"sync"
	"path/filepath"
	"time"
	"io"
	"strconv"
)

func NewFileLogger(filePath string, level, multiSize int) (ILogger, error) {
	os.MkdirAll(filePath, os.ModePerm)
	lf := &LogFile{
		path:      filePath,
		multiSize: multiSize,
	}
	err := lf.Create()
	if err != nil {
		return nil, err
	}
	l := &FileLogger{
		Level: level,
		Out:   lf,
	}
	return l, nil
}

type FileLogger struct {
	Out   io.WriteCloser
	Level int
}

func (l *FileLogger) SetLevel(level int) {
	l.Level = level
}

func (l *FileLogger) Trace(params ...interface{}) {
	if l.Level <= LogTrace {
		l.Out.Write([]byte((fmt.Sprintf("%s [TRACE] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Debug(params ...interface{}) {
	if l.Level <= LogDebug {
		l.Out.Write([]byte((fmt.Sprintf("%s [DEBUG] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Info(params ...interface{}) {
	if l.Level <= LogInfo {
		l.Out.Write([]byte((fmt.Sprintf("%s [INFO] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Error(params ...interface{}) {
	if l.Level <= LogError {
		l.Out.Write([]byte((fmt.Sprintf("%s [ERROR] %s\n", Now(), fmt.Sprint(params ...)))))
	}
}
func (l *FileLogger) Tracef(format string, params ...interface{}) {
	if l.Level <= LogTrace {
		l.Out.Write([]byte((fmt.Sprintf("%s [TRACE] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Debugf(format string, params ...interface{}) {
	if l.Level <= LogDebug {
		l.Out.Write([]byte((fmt.Sprintf("%s [DEBUG] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Infof(format string, params ...interface{}) {
	if l.Level <= LogInfo {
		l.Out.Write([]byte((fmt.Sprintf("%s [INFO] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}
func (l *FileLogger) Errorf(format string, params ...interface{}) {
	if l.Level <= LogError {
		l.Out.Write([]byte((fmt.Sprintf("%s [ERROR] %s\n", Now(), fmt.Sprintf(format, params...)))))
	}
}

func (l *FileLogger) Close() error {
	return l.Out.Close()
}

type LogFile struct {
	path      string
	file      *os.File
	multiSize int
	size      int
	index     int
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
	l.index ++
	file, err := os.OpenFile(filepath.Join(l.path, time.Now().Format("2006-01-02-150405")+"_"+strconv.Itoa(l.index)+".log"),
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
