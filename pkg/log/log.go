package log

import (
	"context"
	stdlog "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogHelper interface {
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
}

// Init initialize a logger with options
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	std = New(opts)
}

// New creat a logger with options
func New(opts *Options) *Logger {
	if opts == nil {
		opts = NewOptions() // fall back to default options
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	//  serializes a Level to an all-caps string
	encodeLevel := zapcore.CapitalLevelEncoder
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder // serializes a Level to an all-caps string and adds color
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time", // key for time in json format
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel, // log level
		EncodeTime:     timeEncoder, // time format
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // trimming all but the final directory from the full path.
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       opts.Development,       // puts the logger in development mode
		DisableCaller:     opts.DisableCaller,     //  stops annotating logs with the calling function's file name and line number
		DisableStacktrace: opts.DisableStacktrace, // completely disables automatic stacktrace capturing
		Sampling: &zap.SamplingConfig{ // nil SamplingConfig disables sampling
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         opts.Format, // sets the logger's encoding
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	// Build() generates a *zap.Logger object based on current configuation
	// zap.AddStacktrace configures the Logger to record a stack trace for all messages at or above a given level.
	// zap.AddCallerSkip  increases the number of callers skipped by caller annotation (as enabled by the AddCaller option).
	l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	logger := &Logger{
		Logger: l,

		// Sometimes we wish to ouput file name and row number where the function's called. In which we use zap.AddCallerSkip(skip int) to skip to the upper level.
		skipCaller: l.WithOptions(zap.AddCallerSkip(1)),

		minLevel:         zapLevel,
		errorStatusLevel: zap.ErrorLevel,
		caller:           true,
		withTraceID:      true,
		//stackTrace:       true,
	}

	/*
		RedirectStdLog redirects output from the standard library's package-global logger to the supplied logger at InfoLevel. Since zap already handles caller annotations, timestamps, etc., it automatically disables the standard library's annotations and prefixing.
	*/
	zap.RedirectStdLog(l)

	return logger
}

// ZapLogger used for other log wrapper such as klog.
func ZapLogger() *zap.Logger {
	return std.Logger
}

// CheckIntLevel used for other log wrapper such as klog which return if logging a
// message at the specified level is enabled.
func CheckIntLevel(level int32) bool {
	var lvl zapcore.Level
	if level < 5 {
		lvl = zapcore.InfoLevel
	} else {
		lvl = zapcore.DebugLevel
	}
	checkEntry := std.Logger.Check(lvl, "")

	return checkEntry != nil
}

// Debug method output debug level log.
func Debug(msg string, fields ...Field) {
	std.Logger.Debug(msg, fields...)
}

// DebugC method output debug level log.
func DebugC(ctx context.Context, msg string, fields ...Field) {
	std.DebugContext(ctx, msg, fields...)
}

// Debugf method output debug level log.
func Debugf(format string, v ...interface{}) {
	std.Logger.Sugar().Debugf(format, v...)
}

// DebugfC method output debug level log.
func DebugfC(ctx context.Context, format string, v ...interface{}) {
	std.DebugfContext(ctx, format, v...)
}

// Debugw method output debug level log.
func Debugw(msg string, keysAndValues ...interface{}) {
	std.Logger.Sugar().Debugw(msg, keysAndValues...)
}

func DebugwC(ctx context.Context, msg string, keysAndValues ...interface{}) {
	std.DebugfContext(ctx, msg, keysAndValues...)
}

// Info method output info level log.
func Info(msg string, fields ...Field) {
	std.Logger.Info(msg, fields...)
}

func InfoC(ctx context.Context, msg string, fields ...Field) {
	std.InfoContext(ctx, msg, fields...)
}

// Infof method output info level log.
func Infof(format string, v ...interface{}) {
	std.Logger.Sugar().Infof(format, v...)
}

func InfofC(ctx context.Context, format string, v ...interface{}) {
	std.InfofContext(ctx, format, v...)
}

// Warn method output warning level log.
func Warn(msg string, fields ...Field) {
	std.Logger.Warn(msg, fields...)
}

func WarnC(ctx context.Context, msg string, fields ...Field) {
	std.WarnContext(ctx, msg, fields...)
}

// Warnf method output warning level log.
func Warnf(format string, v ...interface{}) {
	std.Logger.Sugar().Warnf(format, v...)
}

func WarnfC(ctx context.Context, format string, v ...interface{}) {
	std.WarnfContext(ctx, format, v...)
}

// Error method output error level log.
func Error(msg string, fields ...Field) {
	std.Logger.Error(msg, fields...)
}

func ErrorC(ctx context.Context, msg string, fields ...Field) {
	std.ErrorContext(ctx, msg, fields...)
}

// Errorf method output error level log.
func Errorf(format string, v ...interface{}) {
	std.Logger.Sugar().Errorf(format, v...)
}

func ErrorfC(ctx context.Context, format string, v ...interface{}) {
	std.ErrorfContext(ctx, format, v...)
}

// Panic method output panic level log and shutdown application.
func Panic(msg string, fields ...Field) {
	std.Logger.Panic(msg, fields...)
}

func PanicC(ctx context.Context, msg string, fields ...Field) {
	std.PanicContext(ctx, msg, fields...)
}

// Panicf method output panic level log and shutdown application.
func Panicf(format string, v ...interface{}) {
	std.Logger.Sugar().Panicf(format, v...)
}

func PanicfC(ctx context.Context, format string, v ...interface{}) {
	std.PanicfContext(ctx, format, v...)
}

// Fatal method output fatal level log.
func Fatal(msg string, fields ...Field) {
	std.Logger.Fatal(msg, fields...)
}

func FatalC(ctx context.Context, msg string, fields ...Field) {
	std.PanicContext(ctx, msg, fields...)
}

// Fatalf method output fatal level log.
func Fatalf(format string, v ...interface{}) {
	std.Logger.Sugar().Fatalf(format, v...)
}

func FatalfC(ctx context.Context, format string, v ...interface{}) {
	std.FatalfContext(ctx, format, v...)
}

func StdInfoLogger() *stdlog.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.Logger, zapcore.InfoLevel); err == nil {
		return l
	}

	return nil
}

func Flush() { std.Flush() }
