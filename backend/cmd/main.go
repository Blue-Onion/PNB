package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	aichecker "github.com/Blue-Onion/txShield/internal/AiChecker"
	decoder "github.com/Blue-Onion/txShield/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/txShield/internal/PermissionChecker"
	report "github.com/Blue-Onion/txShield/internal/Report"
	scamchecker "github.com/Blue-Onion/txShield/internal/ScamChecker"
	"github.com/Blue-Onion/txShield/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("Config Not loaded")
	}
	port := cfg.Port
	server := http.Server{
		Handler: newServerRouter(),
		Addr:    ":" + port,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	go func() {
		log.Printf("Listening on http://localhost:%s", port)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error occurred: %v", err)
		}
	}()
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error occurred in Shutdown: %v", err)
	}
}

func newServerRouter() http.Handler {
	router := chi.NewRouter()
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
		w.Write([]byte("hello"))
	})
	return router
}

var errInvalidCalldata = errors.New("invalid calldata")

// checkCode runs the transaction checks and returns the JSON report.
func checkCode(hexCode string) ([]byte, error) {
	res, err := decoder.Decoder(hexCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidCalldata, err)
	}

	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		log.Println(err)
	}
	cfg, err := config.Load()
	if err != nil {
		return []byte(""), err
	}
	ai, err := aichecker.Check(res, cfg.GeminiAPIKey)
	if err != nil {
		log.Println(err)
	}

	reportID := uuid.New().String()
	result := report.Build(reportID, res, ai)
	if err != nil {
		result.Analysis.Status = "AI_UNAVAILABLE"
	}
	return json.MarshalIndent(result, "", "  ")
}
