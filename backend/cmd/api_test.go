package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckCodeAPI(t *testing.T) {
	handler := newAPI("PNB-001", func(hexCode, reportID string) ([]byte, error) {
		if hexCode == "0x00" {
			return checkCode(hexCode, reportID, "")
		}
		if hexCode != "0xdeadbeef" || reportID != "PNB-001" {
			t.Fatalf("unexpected pipeline input %q, %q", hexCode, reportID)
		}
		return []byte(`{"report_id":"PNB-001"}`), nil
	})

	for _, tc := range []struct {
		name, method, contentType, body string
		status                          int
	}{
		{"success", "POST", "application/json", `{"hexcode":" 0xdeadbeef "}`, http.StatusOK},
		{"missing input", "POST", "application/json", `{}`, http.StatusBadRequest},
		{"invalid calldata", "POST", "application/json", `{"hexcode":"0x00"}`, http.StatusBadRequest},
		{"invalid JSON", "POST", "application/json", `{`, http.StatusBadRequest},
		{"unsupported type", "POST", "text/plain", `{}`, http.StatusUnsupportedMediaType},
		{"oversized body", "POST", "application/json", strings.Repeat(" ", (64<<10)+1), http.StatusRequestEntityTooLarge},
		{"wrong method", "GET", "", "", http.StatusMethodNotAllowed},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/api/check-code", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != tc.status {
				t.Fatalf("got HTTP %d, want %d: %s", res.Code, tc.status, res.Body.String())
			}
			if tc.status == http.StatusOK && res.Body.String() != `{"report_id":"PNB-001"}` {
				t.Fatalf("unexpected report: %s", res.Body.String())
			}
		})
	}
}
