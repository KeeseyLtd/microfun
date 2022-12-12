package logging

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var logger *zap.SugaredLogger

func init() {
	var z *zap.Logger
	var err error

	if os.Getenv("APP_ENV") == "testing" {
		wd, _ := os.Getwd()
		for !strings.HasSuffix(wd, "wedding-qr") {
			wd = filepath.Dir(wd)
		}

		loggerCfg := zap.NewDevelopmentConfig()
		loggerCfg.ErrorOutputPaths = []string{fmt.Sprintf("%s/app.log", wd)}
		loggerCfg.OutputPaths = []string{fmt.Sprintf("%s/app.log", wd)}

		z, err = loggerCfg.Build()
	} else {
		loggerCfg := &zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Encoding:         "json",
			EncoderConfig:    encoderConfig,
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		z, err = loggerCfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	}

	defer z.Sync()

	if err != nil {
		z = zap.NewNop()
	}

	logger = z.Sugar()
}

func NewContext(ctx context.Context, args ...interface{}) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(args...))
}

func CopyLoggerContext(newCtx context.Context, oldCtx context.Context) context.Context {
	return context.WithValue(newCtx, loggerKey, oldCtx.Value(loggerKey))
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return logger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return ctxLogger
	}

	return logger
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel(),
	EncodeTime:     zapcore.RFC3339TimeEncoder,
	EncodeDuration: zapcore.MillisDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func encodeLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}
