package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wave-95/boards/server/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestRequestLogger(t *testing.T) {
	t.Run("handler should add duration, request ID, and correlation ID fields", func(t *testing.T) {
		l, observer := logger.NewTest()
		middleware := RequestLogger(l)
		mux := http.NewServeMux()
		handler := middleware(mux)

		rec := httptest.NewRecorder()
		req := buildRequest("", "")

		handler.ServeHTTP(rec, req)
		entries := observer.All()
		//TODO: make test not flaky
		log := entries[0]
		assert.Equal(t, FieldRequestID, log.Context[0].Key)
		assert.Equal(t, FieldCorrelationID, log.Context[1].Key)
		assert.Equal(t, FieldDuration, log.Context[2].Key)
		assert.Equal(t, "[STATUS 404] GET: /", log.Entry.Message)
	})

}
func Test_getOrCreateIDs(t *testing.T) {
	req := buildRequest("", "")
	reqId, corrId := getOrCreateIDs(req)
	assert.NotEqual(t, reqId, "")
	assert.NotEqual(t, corrId, "")
}

func Test_getRequestId(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set(HeaderRequestID, "123abc")
	reqId := getRequestID(req)
	assert.Equal(t, reqId, "123abc")
}

func Test_getCorrelationId(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	req.Header.Set(HeaderCorrelationID, "123abc")
	corrId := getCorrelationID(req)
	assert.Equal(t, corrId, "123abc")
}

func buildRequest(reqId, corrId string) *http.Request {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	if reqId != "" {
		req.Header.Set(HeaderRequestID, reqId)
	}
	if corrId != "" {
		req.Header.Set(HeaderCorrelationID, corrId)
	}
	return req
}
