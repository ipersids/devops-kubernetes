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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("pingpong app: variable PORT is required")
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		log.Fatal("Variable FILE_PATH is required")
	}

	err := os.WriteFile(filePath, []byte("Ping / Pongs: 0"), 0644)
	if err != nil {
		log.Fatalf("Failed create file '%s'", filePath)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var counter atomic.Uint64

	mux := http.NewServeMux()
	mux.HandleFunc("GET /pingpong", func(w http.ResponseWriter, _ *http.Request) {
		current := counter.Add(1) - 1
		status := fmt.Sprintf("Ping / Pongs: %d\n", current)
		err := os.WriteFile(filePath, []byte(string(status)), 0644)
		if err != nil {
			log.Printf("Failed to write file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte(status)); err != nil {
			log.Printf("writing response: %v", err)
		}
		fmt.Print(status)
	})

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
