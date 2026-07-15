package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var counter atomic.Uint64

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("pingpong app: variable PORT is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /pingpong", handleRoot)
	mux.HandleFunc("GET /pings", handlePings)

	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server started in port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}
}

// HANDLERS

func handleRoot(w http.ResponseWriter, r *http.Request) {
	current := counter.Add(1) - 1
	status := fmt.Sprintf("Ping / Pongs: %d\n", current)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(status)); err != nil {
		log.Printf("writing response: %v", err)
	}
	fmt.Printf("%s -> %d\n", r.Pattern, current)
}

func handlePings(w http.ResponseWriter, r *http.Request) {
	current := counter.Load()
	status := fmt.Sprintf("Ping / Pongs: %d\n", current)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(status)); err != nil {
		log.Printf("writing response: %v", err)
	}
	fmt.Printf("%s -> %d\n", r.Pattern, current)
}
