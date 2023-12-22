// Package log 日志框架 zap
package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/woshikedayaa/news-poster/pkg/utils/futil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var (
	filename                               = filepath.Join(futil.GetRootDir(), "logs", "rolling-log.log")
	globalLogger        *zap.Logger        = nil
	globalSugaredLogger *zap.SugaredLogger = nil
	msgKey                                 = "msg"
	levelKey                               = "level"
	timeKey                                = "ts"
	callerKey                              = "caller"
	stacktraceKey                          = "trace"
	globalServiceName                      = "unknown"
)

func New() Logger {
	return newZapLoggerWrapper()
}

func NewSugared() SugaredLogger {
	return newSugaredZapLoggerWrapper()
}

func InitLog(serviceName string) {
	if globalLogger != nil {
		return
	}

	if len(serviceName) != 0 {
		globalServiceName = serviceName
	}

	// encoder
	EncoderConfig := zapcore.EncoderConfig{
		MessageKey:          msgKey,
		LevelKey:            levelKey,
		TimeKey:             timeKey,
		NameKey:             zapcore.OmitKey,
		CallerKey:           callerKey,
		FunctionKey:         zapcore.OmitKey,
		StacktraceKey:       stacktraceKey,
		SkipLineEnding:      false,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         zapcore.CapitalLevelEncoder,
		EncodeTime:          zapcore.RFC3339TimeEncoder,
		EncodeDuration:      nil, // default SecondsDurationEncoder
		EncodeCaller:        nil, // default FullCallerEncoder
		EncodeName:          nil, // default FullNameEncoder
		NewReflectedEncoder: nil, // default defaultReflectedEncoder // private
		ConsoleSeparator:    zapcore.OmitKey,
	}

	encoder := zapcore.NewConsoleEncoder(EncoderConfig)

	//write sync
	logRolling := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    50 * 1024 * 1024,
		MaxAge:     30, // 30 days
		MaxBackups: 30, // 30 (和MaxAge保持一致吧 大于MaxAge的也不会被保存下来)
		LocalTime:  true,
		Compress:   true,
	}

	// cast to buffed
	buffedLoggWS := &zapcore.BufferedWriteSyncer{
		// lumberjack and os.stdout
		WS:            zapcore.NewMultiWriteSyncer(zapcore.AddSync(logRolling), zapcore.AddSync(os.Stdout)),
		Size:          16 * 1024, // 16KB
		FlushInterval: 3,         // 3s
	}

	// enc Encoder, ws WriteSyncer, enab LevelEnabler
	// TODO 把level配置到其他地方
	core := zapcore.NewCore(encoder, buffedLoggWS, zapcore.DebugLevel)
	globalLogger = zap.New(
		core,
		zap.AddCallerSkip(1),
		zap.Fields(zap.String("caller", globalServiceName)),
		zap.AddStacktrace(zapcore.DebugLevel),
	)

	globalSugaredLogger = globalLogger.Sugar()
}

func handleSyncError(logger Logger, err error) {
	logger.Error("error when sync", zap.Error(err))
}
