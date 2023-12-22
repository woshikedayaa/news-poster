package log

import (
	"go.uber.org/zap"
)

// BasicLogger 对应5种不同等级的日志
// 只提供了最基础的日志记录功能
type BasicLogger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

type Logger BasicLogger

type ZapLoggerWrapper struct {
	self *zap.Logger
}

// BasicLogger

func (z *ZapLoggerWrapper) Debug(msg string, fields ...zap.Field) {
	z.self.Debug(msg, fields...)
}

func (z *ZapLoggerWrapper) Info(msg string, fields ...zap.Field) {
	z.self.Info(msg, fields...)
}

func (z *ZapLoggerWrapper) Warn(msg string, fields ...zap.Field) {
	//TODO 上传监控
	z.self.Warn(msg, fields...)
}

func (z *ZapLoggerWrapper) Error(msg string, fields ...zap.Field) {
	//TODO 上传监控
	z.self.Error(msg, fields...)
}

func (z *ZapLoggerWrapper) Fatal(msg string, fields ...zap.Field) {
	//TODO 上传监控
	z.self.Fatal(msg, fields...)
}

func (z *ZapLoggerWrapper) Sync() error {
	return z.self.Sync()
}

func NewZapLoggerWrapper() *ZapLoggerWrapper {
	return &ZapLoggerWrapper{self: globalLogger}
}
