package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewServerRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(logRequests)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Mount("/api", Router())
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("hello")); err != nil {
			requestLogger(r.Context()).Error("Write HTTP response", "error", err)
		}
	})
	return router
}

func Router() *chi.Mux {
	router := chi.NewMux()
	router.Get("/check/{hexcode}", CheckCode())
	return router
}
func CheckCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hexcode := chi.URLParam(r, "hexcode")
		logger := requestLogger(r.Context())
		res, err := checkCode(logger, hexcode)
		if err != nil {
			respondWithJSON(logger, w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(logger, w, http.StatusOK, json.RawMessage(res))
	}
}

type res struct {
	Success bool
	Data    any
}

func respondWithJSON(logger *slog.Logger, w http.ResponseWriter, code int, payload any) {
	res := res{}
	if code < 299 {
		res.Success = true
	} else {
		res.Success = false
	}
	res.Data = payload
	data, err := json.Marshal(res)
	if err != nil {
		logger.Error("Encode JSON response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "Application/Json")
	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		logger.Error("Write JSON response", "error", err)
	}

}
