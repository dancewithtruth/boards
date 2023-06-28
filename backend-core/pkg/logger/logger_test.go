package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
}

func Test_logger_With(t *testing.T) {
	l, observer := NewTestLogger()
	l = l.With("requestID", "123")
	l.Info("Logging with request ID field appended")
	logs := observer.All()
	entry := logs[0]
	fields := entry.Context
	assert.Equal(t, "requestID", fields[0].Key)
	assert.Equal(t, "123", fields[0].String)
}

func Test_logger_WithoutCaller(t *testing.T) {
	l, observer := NewTestLogger()
	l = l.WithoutCaller()
	l.Info("Logging without caller")
	logs := observer.All()
	entry := logs[0]
	assert.Equal(t, false, entry.Caller.Defined)
}

func Test_logger_FromContext(t *testing.T) {
	t.Run("context with logger", func(t *testing.T) {
		l, observer := NewTestLogger()
		l = l.With("customArgKey", "customArgValue")
		ctx := context.WithValue(context.Background(), LoggerKey, l)
		loggerFromCtx := FromContext(ctx)
		assert.NotNil(t, loggerFromCtx)
		loggerFromCtx.Info("Logging to test for customArg field")
		logs := observer.All()
		entry := logs[0]
		fields := entry.Context
		assert.Equal(t, "customArgKey", fields[0].Key)
	})

	t.Run("context without logger", func(t *testing.T) {
		loggerFromCtx := FromContext(context.Background())
		assert.NotNil(t, loggerFromCtx)
	})
}
