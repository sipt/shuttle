package shuttle

import "fmt"

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

type stdLogger struct{}

func (s *stdLogger) Trace(params ...interface{}) {
	fmt.Println("[TRACE]", fmt.Sprint(params ...))
}
func (s *stdLogger) Debug(params ...interface{}) {
	fmt.Println("[DEBUG]", fmt.Sprint(params ...))
}
func (s *stdLogger) Info(params ...interface{}) {
	fmt.Println("[INFO]", fmt.Sprint(params ...))
}
func (s *stdLogger) Error(params ...interface{}) {
	fmt.Println("[ERROR]", fmt.Sprint(params ...))
}
func (s *stdLogger) Tracef(fromat string, params ...interface{}) {
	fmt.Printf("[TRACE] "+fromat+"\n", params...)
}
func (s *stdLogger) Debugf(fromat string, params ...interface{}) {
	fmt.Printf("[DEBUG] "+fromat+"\n", params...)
}
func (s *stdLogger) Infof(fromat string, params ...interface{}) {
	fmt.Printf("[INFO] "+fromat+"\n", params...)
}
func (s *stdLogger) Errorf(fromat string, params ...interface{}) {
	fmt.Printf("[ERROR] "+fromat+"\n", params...)
}
