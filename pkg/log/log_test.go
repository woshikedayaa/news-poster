package log

import "testing"

func TestLog(t *testing.T) {
	InitLog("test")

	logger := New()
	defer logger.Sync()
	logger.Debug("test-debug")
	logger.Info("test-info")
	logger.Warn("test-warn")
	logger.Error("test-error")
	//logger.Fatal("test-fatal")
}

func TestStackTrace(t *testing.T) {
	InitLog("test")
	logger := New()
	defer logger.Sync()
	stTarget(logger)
}

func stTarget(logger Logger) {
	logger.Info("st")
}
