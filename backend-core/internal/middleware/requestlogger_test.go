package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wave-95/boards/backend-core/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	t.Run("handler should add duration, request ID, and correlation ID fields", func(t *testing.T) {
		l, observer := logger.NewTestLogger()
		middleware := RequestLogger(l)
		mux := http.NewServeMux()
		handler := middleware(mux)

		rec := httptest.NewRecorder()
		req := buildRequest("", "")

		handler.ServeHTTP(rec, req)
		entries := observer.All()
		//TODO: make test not flaky
		log := entries[0]
		assert.Equal(t, fieldRequestID, log.Context[0].Key)
		assert.Equal(t, fieldCorrelationID, log.Context[1].Key)
		assert.Equal(t, fieldDuration, log.Context[2].Key)
		assert.Equal(t, "[404] GET: /", log.Entry.Message)
	})
}
func Test_getOrCreateIDs(t *testing.T) {
	req := buildRequest("", "")
	reqID, corrID := getOrCreateIDs(req)
	assert.NotEqual(t, reqID, "")
	assert.NotEqual(t, corrID, "")
}

func Test_getRequestId(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set(headerNameRequestID, "123abc")
	reqID := getRequestID(req)
	assert.Equal(t, reqID, "123abc")
}

func Test_getCorrelationId(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set(headerNameCorrelationID, "123abc")
	corrID := getCorrelationID(req)
	assert.Equal(t, corrID, "123abc")
}

func buildRequest(reqID, corrID string) *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	if reqID != "" {
		req.Header.Set(headerNameRequestID, reqID)
	}
	if corrID != "" {
		req.Header.Set(headerNameCorrelationID, corrID)
	}
	return req
}
