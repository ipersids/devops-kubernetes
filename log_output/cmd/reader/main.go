package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LogoutputHandler struct {
	logsPath     string
	getPingsURL  string
	msgENV       string
	infoFilePath string
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variable PORT is required")
	}

	logsPath := os.Getenv("LOGS_PATH")
	if logsPath == "" {
		log.Fatal("Variable LOGS_PATH is required")
	}

	getPingsURL := os.Getenv("GET_PINGS")
	if getPingsURL == "" {
		log.Fatal("Variable GET_PINGS is required")
	}

	msgENV := os.Getenv("MESSAGE")
	if msgENV == "" {
		log.Fatal("Variable MESSAGE is required")
	}

	infoFilePath := os.Getenv("INFO_FILE_PATH")
	if infoFilePath == "" {
		log.Fatal("Variable INFO_FILE_PATH is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := LogoutputHandler{
		logsPath:     logsPath,
		getPingsURL:  getPingsURL,
		msgENV:       msgENV,
		infoFilePath: infoFilePath,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/logoutput", app.handleLogoutput)

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

func (lh *LogoutputHandler) handleLogoutput(w http.ResponseWriter, r *http.Request) {
	logsData, err := os.ReadFile(lh.logsPath)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	infoData, err := os.ReadFile(lh.infoFilePath)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := http.Get(lh.getPingsURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Printf("Fetching pings failed: error: %v", err)
		http.Error(w, "Pingpong service is not available", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	pongsData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading pings response body failed: error: %v", err)
		http.Error(w, "Somthing went wrong", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("file content: %s\nenv variable: MESSAGE=%s\n%s%s", string(infoData), lh.msgENV, string(logsData), string(pongsData))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("writing response: %v", err)
	}
	fmt.Print(message)
}
