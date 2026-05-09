package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Init initializes the global zap logger. Must be called before any logging.
func Init(env string) (*zap.Logger, error) {
	var err error
	once.Do(func() {
		globalLogger, err = build(env)
		if err == nil {
			zap.ReplaceGlobals(globalLogger)
		}
	})
	return globalLogger, err
}

// L returns the global logger (panics if Init has not been called).
func L() *zap.Logger {
	if globalLogger == nil {
		panic("logger not initialized: call logger.Init first")
	}
	return globalLogger
}

// S returns the global sugared logger.
func S() *zap.SugaredLogger {
	return L().Sugar()
}

func build(env string) (*zap.Logger, error) {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zapcore.Level
	if env == "production" {
		level = zapcore.InfoLevel
	} else {
		level = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)

	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	return zap.New(core, opts...), nil
}
