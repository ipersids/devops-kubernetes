package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// SQL

const createTodosTable = `
CREATE TABLE IF NOT EXISTS todos (
	id UUID DEFAULT uuidv4() PRIMARY KEY,
	title TEXT NOT NULL,
	completed BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ DEFAULT now()
);
INSERT INTO todos (title, completed)
SELECT *
FROM (
    VALUES
        ('Make the project respond something to a GET request sent to the / url', true),
        ('Use a NodePort Service to enable access', true),
        ('External access with Ingress', false),
        ('Develop a second application that simply responds with "pong 0" to a GET request', false)
) AS initial_todos(title, completed)
WHERE NOT EXISTS (SELECT 1 FROM todos);
`

const selectTasks = `
SELECT id, title, completed FROM todos;
`

const insertTask = `
INSERT INTO todos (title)
VALUES ($1)
`

// end SQL

type store struct {
	db *sql.DB
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
		log.Fatal().Msg("Variable PORT is required")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal().Msg("pingpong app: variable DB_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	_, err = db.Exec(createTodosTable)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	app := store{db: db}

	log.Info().Msg("Connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", app.handleTodos)
	mux.HandleFunc("POST /todos", app.handleCreateTodo)

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
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	<-ctx.Done()
	log.Print("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}
}

func (s *store) handleTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := s.getTodos()
	if err != nil {
		logger(r, http.StatusInternalServerError).Err(err).Msg("failed to get todos from database")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(todoResponse{Data: todos})

	if err != nil {
		log.Error().Err(err).Msg("failed encode response")
		return
	}

	logger(r, http.StatusOK).Msg("successful request")
}

func (s *store) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger(r, http.StatusBadRequest).Err(err).Msg("failed to parse form")
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if title == "" || len(title) > 140 {
		logger(r, http.StatusBadRequest).Msg("title should be between 1 and 140 bytes")
		http.Error(w, "Title should be between 1 and 140 bytes", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(insertTask, title)
	if err != nil {
		logger(r, http.StatusInternalServerError).Err(err).Msg("failed to add a new task to database")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger(r, http.StatusCreated).Msg("successful request")
}

// Database

func (s *store) getTodos() ([]todo, error) {
	rows, err := s.db.Query(selectTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []todo

	for rows.Next() {
		var todo todo

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// Logging

func logger(r *http.Request, status int) *zerolog.Event {
	logger := log.Info()
	if status >= 400 {
		logger = log.Error()
	}

	return logger.
		Str("method", r.Method).
		Str("host", r.Host).
		Str("path", r.URL.Path).
		Int("status", status).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent())
}
