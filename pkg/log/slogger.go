package log

import "go.uber.org/zap"

type SugaredLogger interface {
	BasicLogger
	DebugF(f string, args ...interface{})
	InfoF(f string, args ...interface{})
	WarnF(f string, args ...interface{})
	ErrorF(f string, args ...interface{})
	FatalF(f string, args ...interface{})
}

type SugaredZapLoggerWrapper struct {
	*ZapLoggerWrapper
	sugar *zap.SugaredLogger
}

func NewSugaredZapLoggerWrapper() SugaredLogger {
	return &SugaredZapLoggerWrapper{
		ZapLoggerWrapper: NewZapLoggerWrapper(),
		sugar:            globalSugaredLogger,
	}
}

// SugaredLogger

func (z *ZapLoggerWrapper) DebugF(f string, args ...interface{}) {
	z.self.Sugar().Debugf(f, args...)
}

func (z *ZapLoggerWrapper) InfoF(f string, args ...interface{}) {
	z.self.Sugar().Infof(f, args...)
}

func (z *ZapLoggerWrapper) WarnF(f string, args ...interface{}) {
	z.self.Sugar().Warnf(f, args...)
}

func (z *ZapLoggerWrapper) ErrorF(f string, args ...interface{}) {
	z.self.Sugar().Errorf(f, args...)
}

func (z *ZapLoggerWrapper) FatalF(f string, args ...interface{}) {
	z.self.Sugar().Fatalf(f, args...)
}
