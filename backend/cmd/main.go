package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Blue-Onion/txShield/internal/api"
	"github.com/Blue-Onion/txShield/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Load configuration: %v", err)
	}
	port := cfg.Port
	server := http.Server{
		Handler: api.NewServerRouter(),
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
