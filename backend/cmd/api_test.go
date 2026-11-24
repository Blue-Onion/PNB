package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckCodeAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/check/0x00", nil)
	recorder := httptest.NewRecorder()
	newServerRouter().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got HTTP %d, want %d: %s", recorder.Code, http.StatusBadRequest, recorder.Body.String())
	}

	var response struct {
		Success bool
		Data    string
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if response.Success || !strings.Contains(response.Data, "invalid calldata") {
		t.Fatalf("unexpected error response: %s", recorder.Body.String())
	}
}
