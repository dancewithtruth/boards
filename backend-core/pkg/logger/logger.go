package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type keyLogger int

// LoggerKey is the key used to retrieve the logger value from the request context.
const LoggerKey keyLogger = 0

// Logger is an interface that describes all the capabilities of the application's logger.
type Logger interface {
	With(args ...interface{}) Logger

	WithoutCaller() Logger
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})
	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
	// Fatalf uses fmt.Sprintf to construct and log a message, then calls os.Exit.
	Fatalf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewSugar(l)
}

// NewTestLogger returns a test logger with observability capabilities.
func NewTestLogger() (Logger, *observer.ObservedLogs) {
	observedLogs, logObserver := observer.New(zap.InfoLevel)
	testLogger := zap.New(observedLogs)
	return NewSugar(testLogger), logObserver
}

// NewSugar returns a SugaredLogger and implements the Logger interface.
func NewSugar(l *zap.Logger) Logger {
	return &logger{l.Sugar()}
}

// With returns a logger based off the root logger and decorates it with the arguments.
func (l *logger) With(args ...interface{}) Logger {
	return &logger{l.SugaredLogger.With(args...)}
}

// WithoutCaller returns a logger that does not output the caller field and location of the calling code.
func (l *logger) WithoutCaller() Logger {
	return &logger{l.SugaredLogger.WithOptions(zap.WithCaller(false))}
}

// FromContext returns a logger from context. If none found, instantiate a new logger.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(LoggerKey).(Logger); ok {
		return l
	}
	return New()
}
