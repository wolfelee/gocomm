package jlog

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// JLogger default logger
var JLogger *Logger

// Auto ...
func Auto(err error) Func {
	if err != nil {
		return JLogger.With(zap.Any("err", err.Error())).Error
	}

	return JLogger.Info
}

// Info ...
func Info(msg string, fields ...Field) {
	JLogger.desugar.Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	JLogger.desugar.Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	JLogger.desugar.Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	JLogger.desugar.Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	JLogger.desugar.Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	JLogger.desugar.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	JLogger.desugar.Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	if JLogger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	JLogger.sugar.Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	JLogger.sugar.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	JLogger.sugar.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	JLogger.sugar.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	JLogger.sugar.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	JLogger.sugar.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	JLogger.sugar.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	JLogger.sugar.Fatalf(msg, args...)
}

// Log ...
func (fn Func) Log(msg string, fields ...Field) {
	fn(msg, fields...)
}

// With ...
func With(fields ...Field) *Logger {
	return JLogger.With(fields...)
}

type Context = *_context

type _context struct {
	logger *Logger
	fields []zapcore.Field
}

func (c *_context) Logger() *Logger {
	c.logger.desugar.With(c.fields...)
	return c.logger
}

func (c *_context) AddFields(fields ...zapcore.Field) {
	c.logger = c.logger.With(fields...)
	c.fields = append(c.fields, fields...)
}

type _marker struct{}

var (
	_key = _marker{}
)

// ToContext place *zap.Logger to context
func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, _key, &_context{logger: logger})
}

// Extract
func Extract(ctx context.Context) *_context {
	if ctx == nil {
		return nil
	}
	if _ctx, ok := ctx.Value(_key).(*_context); ok {
		return _ctx
	}
	return nil
}

// T return logger from context or default (trace)
func T(ctx context.Context) *Logger {
	_ctx := Extract(ctx)
	if _ctx != nil {
		return _ctx.Logger()
	}
	return JLogger
}
