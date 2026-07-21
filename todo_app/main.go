package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed templates statics
var files embed.FS

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().Msg("Variable PORT is required")
	}

	volumeDir := os.Getenv("VOLUME_DIR")
	if volumeDir == "" {
		log.Fatal().Msg("Variable VOLUME_DIR is required")
	}

	todoBackendURL := os.Getenv("TODO_BACKEND_URL")
	imageName := "image.jpeg"
	imagePath := volumeDir + imageName

	_, err := os.Stat(imagePath)
	if errors.Is(err, os.ErrNotExist) {
		err = downloadAndReplaceImage(volumeDir, imageName)
	}

	if err != nil {
		log.Fatal().Msg("Failed download image on start")
	}

	tmpl := template.Must(template.ParseFS(files, "templates/index.tmpl"))

	app := todoAppHandler{
		volumeDir:      volumeDir,
		imageName:      imageName,
		tmpl:           tmpl,
		mu:             sync.Mutex{},
		refreshing:     false,
		todoBackendURL: todoBackendURL,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.Handle("GET /statics/", http.FileServerFS(files))
	mux.HandleFunc("GET /", app.handleRoot)
	mux.HandleFunc("POST /task", app.handleCreateTask)
	mux.HandleFunc("GET /image", app.handleImage)

	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Msgf("Server started in port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Server failed to start", "error", err)
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msgf("Server forced to shutdown", "error", err)
	}
}
