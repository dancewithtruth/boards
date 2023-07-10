package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlePingCheck(t *testing.T) {
	// Setup request and response recorder
	req, err := http.NewRequest(http.MethodGet, "localhost:8080/ping", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	// Test ping check
	handlePingCheck(rr, req)
	statusCode := rr.Result().StatusCode
	if statusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %v", statusCode)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "pong") {
		t.Errorf("Expected response to include pong, got: %v", body)
	}
}
