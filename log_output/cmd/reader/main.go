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

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		log.Fatal("Variable FILE_PATH is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/logoutput", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if _, err := w.Write(data); err != nil {
			log.Printf("writing response: %v", err)
		}
		fmt.Printf("%s %s : %s", r.Method, r.Host+r.Pattern, string(data))
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
