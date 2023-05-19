package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/google/uuid"
)

type requestIdKey int
type correlationIdKey int

const (
	RequestIdKey     requestIdKey     = 0
	CorrelationIdKey correlationIdKey = 0

	HeaderRequestID     = "X-Request-ID"
	HeaderCorrelationID = "X-Correlation-ID"

	FieldDuration      = "duration"
	FieldBytes         = "bytes"
	FieldRequestID     = "requestID"
	FieldCorrelationID = "correlationID"
)

// LoggerRW is a wrapper around ResponseWriter meant to capture the status code and number of bytes written
type LoggerRW struct {
	http.ResponseWriter
	StatusCode   int
	BytesWritten int
}

func (lrw *LoggerRW) Write(p []byte) (int, error) {
	bytesWritten, err := lrw.ResponseWriter.Write(p)
	lrw.BytesWritten = bytesWritten
	return bytesWritten, err
}

func (lrw *LoggerRW) WriteHeader(statusCode int) {
	lrw.StatusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// RequestLogger is a middleware that creates a context-aware logger for every request and makes it available to downstream
// handlers on the request context.
func RequestLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Get request and correlation IDs and append fields to request logger
			// Then set logger to request context
			reqId, corrId := getOrCreateIDs(r)
			requestLogger := l.With(FieldRequestID, reqId, FieldCorrelationID, corrId)
			ctx := r.Context()
			ctx = context.WithValue(ctx, logger.LoggerKey, requestLogger)
			r = r.WithContext(ctx)

			// Wrap rw to make status code and bytes written available
			lrw := &LoggerRW{ResponseWriter: rw, StatusCode: http.StatusOK}
			next.ServeHTTP(lrw, r)

			// Log duration of request
			requestLogger.
				WithoutCaller().
				With(FieldDuration, time.Since(start).Milliseconds(), FieldBytes, lrw.BytesWritten).
				Infof("[%v] %s: %s", lrw.StatusCode, r.Method, r.URL.Path)
		})
	}
}

// Look for existing request ID and correlation ID from request header
// Return existing values or generate new uuids if not found
func getOrCreateIDs(r *http.Request) (reqId string, corrId string) {
	reqId = getRequestID(r)
	corrId = getCorrelationID(r)
	if reqId == "" {
		reqId = uuid.NewString()
	}
	if corrId == "" {
		corrId = uuid.NewString()
	}
	return reqId, corrId
}

// getRequestID grabs the request ID string off the X-Request-ID header
func getRequestID(r *http.Request) string {
	return r.Header.Get(HeaderRequestID)
}

// getCorrelationId grabs the correlation ID string off the X-Correlation-ID header
// The correlation id groups together multiple request ids
func getCorrelationID(r *http.Request) string {
	return r.Header.Get(HeaderCorrelationID)
}
