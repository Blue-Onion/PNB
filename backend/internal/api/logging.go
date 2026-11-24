package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type requestLoggerKey struct{}

func requestLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(requestLoggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := slog.With("request_id", uuid.NewString())
		wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		completed := false
		defer func() {
			route := chi.RouteContext(r.Context()).RoutePattern()
			status := wrapped.Status()
			level := slog.LevelInfo
			switch {
			case !completed, status >= http.StatusInternalServerError:
				level = slog.LevelError
			case status >= http.StatusBadRequest:
				level = slog.LevelWarn
			case status == 0:
				status = http.StatusOK
			}
			logger.Log(r.Context(), level, "HTTP request",
				"method", r.Method, "route", route, "status", status,
				"duration", time.Since(start), "completed", completed)
		}()
		next.ServeHTTP(wrapped, r.WithContext(context.WithValue(r.Context(), requestLoggerKey{}, logger)))
		completed = true
	})
}
