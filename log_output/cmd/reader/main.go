package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variable PORT is required")
	}

	logsPath := os.Getenv("LOGS_PATH")
	if logsPath == "" {
		log.Fatal("Variable LOGS_PATH is required")
	}

	pongsPath := os.Getenv("PONGS_PATH")
	if pongsPath == "" {
		log.Fatal("Variable PONGS_PATH is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/logoutput", func(w http.ResponseWriter, r *http.Request) {
		logsData, err := os.ReadFile(logsPath)
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pongsData, err := os.ReadFile(pongsPath)
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		msg := fmt.Sprintf("%s%s", string(logsData), string(pongsData))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Printf("writing response: %v", err)
		}
		fmt.Print(msg)
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
