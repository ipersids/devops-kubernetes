package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

//go:embed templates statics
var files embed.FS

type todo struct {
	ID        string
	Title     string
	Completed bool
}

type pageData struct {
	Title string
	Todos []todo
}

var todoList = []todo{
	{ID: "0", Title: "Make the project respond something to a GET request sent to the / url", Completed: true},
	{ID: "1", Title: "Use a NodePort Service to enable access"},
	{ID: "2", Title: "External access with Ingress"},
	{ID: "3", Title: "Develop a second application that simply responds with \"pong 0\" to a GET request"},
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tmpl := template.Must(template.ParseFS(files, "templates/index.tmpl"))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.Handle("GET /statics/", http.FileServerFS(files))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		data := pageData{
			Title: "My day",
			Todos: todoList,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("render page: %v", err)
		}
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
