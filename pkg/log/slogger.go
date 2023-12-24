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

func newSugaredZapLoggerWrapper() SugaredLogger {
	return &SugaredZapLoggerWrapper{
		ZapLoggerWrapper: newZapLoggerWrapper(),
		sugar:            globalSugaredLogger,
	}
}

// SugaredLogger

func (z *SugaredZapLoggerWrapper) DebugF(f string, args ...interface{}) {
	z.sugar.Debugf(f, args...)
}

func (z *SugaredZapLoggerWrapper) InfoF(f string, args ...interface{}) {
	z.sugar.Infof(f, args...)
}

func (z *SugaredZapLoggerWrapper) WarnF(f string, args ...interface{}) {
	z.sugar.Warnf(f, args...)
}

func (z *SugaredZapLoggerWrapper) ErrorF(f string, args ...interface{}) {
	z.sugar.Errorf(f, args...)
}

func (z *SugaredZapLoggerWrapper) FatalF(f string, args ...interface{}) {
	z.sugar.Fatalf(f, args...)
}

func (z *SugaredZapLoggerWrapper) Sync() {
	//TODO 上传监控
	var err error
	err = z.sugar.Sync()
	handleSyncError(z, err)
	z.ZapLoggerWrapper.Sync()
}
