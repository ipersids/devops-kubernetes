package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
)

var todoList = []todo{
	{
		ID:        uuid.MustParse("49f8d30d-f5de-482e-b67b-fd2e61524080"),
		Title:     "Make the project respond something to a GET request sent to the / url",
		Completed: true,
	},
	{
		ID:        uuid.MustParse("c159ae3a-64c6-4f3b-b800-c62276852f6c"),
		Title:     "Use a NodePort Service to enable access",
		Completed: true,
	},
	{
		ID:        uuid.MustParse("d8b0a007-2aa3-4406-9410-06af8e08fbd0"),
		Title:     "External access with Ingress",
		Completed: false,
	},
	{
		ID:        uuid.MustParse("da8faa02-2a1e-4b2d-8ad2-c2649adcdea6"),
		Title:     "Develop a second application that simply responds with \"pong 0\" to a GET request",
		Completed: false,
	},
}

type todo struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}

type todoRequest struct {
	Title string `json:"title"`
}

type todoResponse struct {
	Data []todo `json:"data"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Variable PORT is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(todoResponse{Data: todoList})

		if err != nil {
			log.Printf("Failed encode response: error: %v", err)
		}
	})

	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("New task: failed to parse form: %v", err)
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")

		if title == "" || len(title) > 140 {
			http.Error(w, "Title should be between 1 and 140 bytes", http.StatusBadRequest)
			return
		}

		newTask := todo{
			ID:        uuid.New(),
			Title:     title,
			Completed: false,
		}

		todoList = append(todoList, newTask)

		http.Redirect(w, r, "/", http.StatusSeeOther)
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
