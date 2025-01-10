package log

import (
	"log"
	"strings"

	"github.com/xakepp35/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EnvLoggerDisable  = "LOGGER_DISABLE"
	EnvLoggerOut      = "LOGGER_OUT"
	EnvLoggerErr      = "LOGGER_ERR"
	EnvLoggerLevel    = "LOGGER_LEVEL"
	EnvLoggerDev      = "LOGGER_DEV"
	EnvLoggerEncoding = "LOGGER_ENCODING"
	EnvLoggerStrace   = "LOGGER_STRACE"
	EnvLoggerColor    = "LOGGER_COLOR"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
	StdLog *log.Logger
)

var atomicLevel zap.AtomicLevel

func init() {
	var res *zap.Logger
	// always set globals
	defer func() {
		Logger = res
		Sugar = res.Sugar()
		zap.ReplaceGlobals(res)
		StdLog = zap.NewStdLog(res)
	}()
	res, err := NewZap()
	if err != nil {
		panic(err)
	}
}

// Initialize the logger.
func NewZap() (res *zap.Logger, err error) {
	// check if completely disable the logger
	loggerDisable := env.Get(EnvLoggerDisable, false)
	if loggerDisable {
		// Use the no-op logger that discards all log messages
		Logger = zap.NewNop()
		return
	}
	// load env
	loggerOut := env.Get(EnvLoggerOut, "stdout")
	loggerErr := env.Get(EnvLoggerErr, "stderr")
	loggerLevel := env.Get(EnvLoggerLevel, "info")
	loggerDev := env.Get(EnvLoggerDev, false)
	loggerEncoding := env.Get(EnvLoggerEncoding, "json")
	loggerStrace := env.Get(EnvLoggerStrace, false)
	loggerColor := env.Get(EnvLoggerColor, false)
	// Custom configuration for Zap
	atomicLevel, _ = zap.ParseAtomicLevel(strings.ToLower(loggerLevel))
	encodeLevel := zapcore.CapitalLevelEncoder
	if loggerColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}
	// config := zap.NewProductionConfig()
	zapConfig := zap.Config{
		Level:       atomicLevel,
		Development: loggerDev,
		Encoding:    loggerEncoding,
		// DisableCaller:  generic.GetEnv(log.EnvLoggereCaller, false),
		DisableStacktrace: !loggerStrace,
		// Sampling: &zap.SamplingConfig{
		// 	Initial:    1000,
		// 	Thereafter: 1000,
		// },
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp", // t
			LevelKey:       "level",     // l
			NameKey:        "zap",       // z
			CallerKey:      "caller",    // c
			MessageKey:     "message",   // m
			StacktraceKey:  "strace",    // s
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{loggerOut},
		ErrorOutputPaths: []string{loggerErr},
	}
	res, err = zapConfig.Build(zap.WithCaller(true), zap.AddCaller())
	if err != nil {
		panic("zapConfig.Build(): " + err.Error())
	}
	return
}

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)
