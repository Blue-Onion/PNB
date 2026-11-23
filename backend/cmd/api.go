package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"strings"
	"time"
)

func serveAPI(addr, cert, key, reportID string, check func(string, string) ([]byte, error)) error {
	if (cert == "") != (key == "") {
		return errors.New("provide both -tls-cert and -tls-key to enable HTTPS")
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           newAPI(reportID, check),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       time.Minute,
	}
	if cert != "" {
		log.Printf("API listening on https://%s", addr)
		return server.ListenAndServeTLS(cert, key)
	}

	log.Printf("API listening on http://%s", addr)
	return server.ListenAndServe()
}

func newAPI(reportID string, check func(string, string) ([]byte, error)) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/check-code", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			writeAPIError(w, http.StatusMethodNotAllowed, "Use POST to check calldata.")
			return
		}

		contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || contentType != "application/json" {
			writeAPIError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json.")
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				writeAPIError(w, http.StatusRequestEntityTooLarge, "Request body exceeds 64 KiB.")
			} else {
				writeAPIError(w, http.StatusBadRequest, "Could not read the request body.")
			}
			return
		}

		var request struct {
			HexCode string `json:"hexcode"`
		}
		if err := json.Unmarshal(body, &request); err != nil {
			writeAPIError(w, http.StatusBadRequest, "Send a JSON object with a hexcode string.")
			return
		}
		hexCode := strings.TrimSpace(request.HexCode)
		if hexCode == "" {
			writeAPIError(w, http.StatusBadRequest, "hexcode is required.")
			return
		}

		data, err := check(hexCode, reportID)
		switch {
		case errors.Is(err, errInvalidCalldata):
			writeAPIError(w, http.StatusBadRequest, err.Error())
		case err != nil:
			log.Printf("build report: %v", err)
			writeAPIError(w, http.StatusInternalServerError, "Could not generate the report.")
		default:
			writeAPIResponse(w, http.StatusOK, apiResponse{Result: data})
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "API route not found.")
	})
	return mux
}

type apiResponse struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeAPIResponse(w, status, apiResponse{Error: message})
}

func writeAPIResponse(w http.ResponseWriter, status int, response apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("write API response: %v", err)
	}
}
