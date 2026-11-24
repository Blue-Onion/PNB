package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTPLogsKeepRequestDataPrivate(t *testing.T) {
	var logs bytes.Buffer
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&logs, nil)))
	t.Cleanup(func() { slog.SetDefault(previousLogger) })

	for _, tc := range []struct {
		method string
		target string
		status int
		level  string
	}{
		{http.MethodGet, "/api/check/private-calldata?token=private-token", http.StatusBadRequest, "WARN"},
		{http.MethodGet, "/private-unmatched", http.StatusNotFound, "WARN"},
		{http.MethodOptions, "/api/check/private-calldata", http.StatusOK, "INFO"},
	} {
		logs.Reset()
		req := httptest.NewRequest(tc.method, tc.target, nil)
		if tc.method == http.MethodOptions {
			req.Header.Set("Origin", "https://example.com")
			req.Header.Set("Access-Control-Request-Method", http.MethodGet)
		}
		recorder := httptest.NewRecorder()
		NewServerRouter().ServeHTTP(recorder, req)
		if recorder.Code != tc.status {
			t.Fatalf("%s %s: got status %d, want %d", tc.method, tc.target, recorder.Code, tc.status)
		}
		if strings.Contains(logs.String(), "private-") {
			t.Fatalf("request data leaked into logs: %s", &logs)
		}
		var entry struct {
			Message   string `json:"msg"`
			Level     string `json:"level"`
			Status    int    `json:"status"`
			Route     string `json:"route"`
			RequestID string `json:"request_id"`
		}
		entries := json.NewDecoder(&logs)
		var errorRequestID string
		for entries.More() {
			if err := entries.Decode(&entry); err != nil {
				t.Fatal(err)
			}
			if entry.Message == "Decode calldata" {
				errorRequestID = entry.RequestID
			}
		}
		if entry.Message != "HTTP request" || entry.Level != tc.level || entry.Status != tc.status || entry.RequestID == "" {
			t.Fatalf("unexpected HTTP request log: %+v", entry)
		}
		if tc.status == http.StatusBadRequest && (entry.Route != "/api/check/{hexcode}" || errorRequestID != entry.RequestID) {
			t.Fatalf("expected route template and matching error request ID: %+v", entry)
		}
	}
}

func TestCheckCodeAPI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/check/0x00", nil)
	recorder := httptest.NewRecorder()
	NewServerRouter().ServeHTTP(recorder, req)

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
