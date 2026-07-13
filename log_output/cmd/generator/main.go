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

	"github.com/google/uuid"
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

	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:         addr,
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

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			id := uuid.NewString()
			now := time.Now().UTC().Format(time.RFC3339Nano)
			data := fmt.Sprintf("%s: %s\n", now, id)
			err := os.WriteFile(filePath, []byte(data), 0644)
			if err != nil {
				log.Fatal("Faled write data "+data, "FILE_PATH=", filePath)
			}
			fmt.Print(data)
			<-ticker.C
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
