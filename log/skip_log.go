package log

func NewSkipLogger() (ILogger, error) {
	return &SkipLogger{}, nil
}

type SkipLogger struct{}

func (s *SkipLogger) SetLevel(level int)                          {}
func (s *SkipLogger) Trace(params ...interface{})                 {}
func (s *SkipLogger) Debug(params ...interface{})                 {}
func (s *SkipLogger) Info(params ...interface{})                  {}
func (s *SkipLogger) Error(params ...interface{})                 {}
func (s *SkipLogger) Tracef(fromat string, params ...interface{}) {}
func (s *SkipLogger) Debugf(fromat string, params ...interface{}) {}
func (s *SkipLogger) Infof(fromat string, params ...interface{})  {}
func (s *SkipLogger) Errorf(fromat string, params ...interface{}) {}
func (s *SkipLogger) Close() error                                { return nil }
