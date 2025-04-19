package log

import (
	"sync"
	"time"

	krtlog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormlogger "gorm.io/gorm/logger"
)

type Field = zapcore.Field

type Logger interface {
	Debugf(format string, args ...any)
	Debugw(msg string, keyvals ...any)
	Infof(format string, args ...any)
	Infow(msg string, keyvals ...any)
	Warnf(format string, args ...any)
	Warnw(msg string, keyvals ...any)
	Errorf(format string, args ...any)
	Errorw(err error, msg string, keyvals ...any)
	Panicf(format string, args ...any)
	Panicw(msg string, keyvals ...any)
	Fatalf(format string, args ...any)
	Fatalw(msg string, keyvals ...any)
	With(fields ...Field) Logger
	AddCallerSkip(skip int) Logger
	Sync()

	// integrate other loggers
	krtlog.Logger
	gormlogger.Interface
}

// zapLogger is implementation of Logger interface.
type zapLogger struct {
	z    *zap.Logger
	opts *Options
}

var _ Logger = (*zapLogger)(nil)

var (
	mu sync.Mutex

	// std is global default logger instance.
	std = NewLogger(NewOptions())
)

func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()

	std = NewLogger(opts)
}

func NewLogger(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	// set zap log level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	// Specify the time serialization function and make the time serialization
	// in `2006-01-02 15:04:05.000` format, which is easier to read.
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// Specify the time.Duration serialization function, and serialize the time.Duration into a floating point number of milliseconds passed.
	// Milliseconds are more accurate than the default seconds
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}
	// when output to local path, with color is forbidden
	if opts.Format == "console" && opts.EnableColor {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// make zap logger config
	cfg := &zap.Config{
		// 是否在日志中显示调用日志所在的文件和行号，例如：`"caller":"onex/onex.go:75"`
		DisableCaller: opts.DisableCaller,
		// 是否禁止在 panic 及以上级别打印堆栈信息
		DisableStacktrace: opts.DisableStacktrace,
		// 指定日志级别
		Level: zap.NewAtomicLevelAt(zapLevel),
		// 指定日志显示格式，可选值：console, json
		Encoding:      opts.Format,
		EncoderConfig: encoderConfig,
		// 制定日志输出位置
		OutputPaths: opts.OutputPaths,
		// 设置 zap 内部错误输出位置
		ErrorOutputPaths: []string{"stderr"},
	}
	z, err := cfg.Build(
		zap.AddStacktrace(zapcore.PanicLevel),
		zap.AddCallerSkip(2),
	)
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{
		z:    z,
		opts: opts,
	}

	// 把标准库的 log.Logger 的 info 级别的输出重定向到 zap.Logger
	zap.RedirectStdLog(z)
	return logger
}

func Default() Logger {
	return std
}

// Sync 调用底层 zap.Logger 的 Sync 方法，将缓存中的日志刷新到磁盘文件中. 主程序需要在退出前调用 Sync.
func Sync() {
	std.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.z.Sugar().Debugf(format, args...)
}

// Debugw 输出 debug 级别的日志.
func Debugw(msg string, keyvals ...any) {
	std.Debugw(msg, keyvals...)
}

func (l *zapLogger) Debugw(msg string, keyvals ...any) {
	l.z.Sugar().Debugw(msg, keyvals...)
}

// Infof 输出 info 级别的日志.
func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

func (l *zapLogger) Infof(msg string, keyvals ...any) {
	l.z.Sugar().Infof(msg, keyvals...)
}

// Infow 输出 info 级别的日志.
func Infow(msg string, keyvals ...any) {
	std.Infow(msg, keyvals...)
}

func (l *zapLogger) Infow(msg string, keyvals ...any) {
	l.z.Sugar().Infow(msg, keyvals...)
}

// Warnf 输出 warning 级别的日志.
func Warnf(format string, args ...any) {
	std.Warnf(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.z.Sugar().Warnf(format, args...)
}

// Warnw 输出 warning 级别的日志.
func Warnw(msg string, keyvals ...any) {
	std.Warnw(msg, keyvals...)
}

func (l *zapLogger) Warnw(msg string, keyvals ...any) {
	l.z.Sugar().Warnw(msg, keyvals...)
}

// Errorf 输出 error 级别的日志.
func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.z.Sugar().Errorf(format, args...)
}

// Errorw 输出 error 级别的日志.
func Errorw(err error, msg string, keyvals ...any) {
	std.Errorw(err, msg, keyvals...)
}

func (l *zapLogger) Errorw(err error, msg string, keyvals ...any) {
	l.z.Sugar().Errorw(msg, append(keyvals, "err", err)...)
}

// Panicf 输出 panic 级别的日志.
func Panicf(format string, args ...any) {
	std.Panicf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...any) {
	l.z.Sugar().Panicf(format, args...)
}

// Panicw 输出 panic 级别的日志.
func Panicw(msg string, keyvals ...any) {
	std.Panicw(msg, keyvals...)
}

func (l *zapLogger) Panicw(msg string, keyvals ...any) {
	l.z.Sugar().Panicw(msg, keyvals...)
}

// Fatalf 输出 fatal 级别的日志.
func Fatalf(format string, args ...any) {
	std.Fatalf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...any) {
	l.z.Sugar().Fatalf(format, args...)
}

// Fatalw 输出 fatal 级别的日志.
func Fatalw(msg string, keyvals ...any) {
	std.Fatalw(msg, keyvals...)
}

func (l *zapLogger) Fatalw(msg string, keyvals ...any) {
	l.z.Sugar().Fatalw(msg, keyvals...)
}

func With(fields ...Field) Logger {
	return std.With(fields...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *zapLogger) With(fields ...Field) Logger {
	if len(fields) == 0 {
		return l
	}

	lc := l.clone()
	lc.z = lc.z.With(fields...)
	return lc
}

// clone 深度拷贝 zapLogger
func (l *zapLogger) clone() *zapLogger {
	copied := *l
	return &copied
}

// AddCallerSkip increases the number of callers skipped by caller annotation
// (as enabled by the AddCaller option). When building wrappers around the
// Logger and SugaredLogger, supplying this Option prevents zap from always
// reporting the wrapper code as the caller.
func (l *zapLogger) AddCallerSkip(skip int) Logger {
	lc := l.clone()
	lc.z = lc.z.WithOptions(zap.AddCallerSkip(skip))
	return lc
}
