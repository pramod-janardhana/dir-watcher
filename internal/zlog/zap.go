package zlog

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	zap.SugaredLogger
	lj *lumberjack.Logger
}

func NewZapLogger(out io.Writer) Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.FunctionKey = "F"
	encoderConfig.CallerKey = "C"

	syslogCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(out),
		zapcore.DebugLevel,
	)

	// Create a new logger with the Syslog core
	log := zap.New(syslogCore,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
	)

	return &ZapLogger{SugaredLogger: *log.Sugar()}
}

func (z *ZapLogger) Close() {
	z.Close()
}
