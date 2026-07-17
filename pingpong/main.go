package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// SQL

const createPingsTable = `
CREATE TABLE IF NOT EXISTS pings (
	id INTEGER PRIMARY KEY,
  amount INTEGER NOT NULL
);
INSERT INTO pings (id, amount)
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;
`

const updatePingsCounter = `
UPDATE pings
SET amount = amount + 1
WHERE id = 1
RETURNING amount;
`

const selectPingsCounter = `
SELECT amount FROM pings;
`

// end SQL

type store struct {
	db *sql.DB
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("pingpong app: variable PORT is required")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("pingpong app: variable DB_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPingsTable)
	if err != nil {
		log.Fatal(err)
	}

	app := store{db: db}

	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /pingpong", app.handleRoot)
	mux.HandleFunc("GET /pings", app.handlePings)

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

func (s *store) handleRoot(w http.ResponseWriter, r *http.Request) {
	var count int

	err := s.db.QueryRow(updatePingsCounter).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	status := fmt.Sprintf("Ping / Pongs: %d\n", count)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(status)); err != nil {
		log.Printf("writing response: %v", err)
	}
	fmt.Printf("%s -> %d\n", r.Pattern, count)
}

func (s *store) handlePings(w http.ResponseWriter, r *http.Request) {
	var count int

	err := s.db.QueryRow(selectPingsCounter).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	status := fmt.Sprintf("Ping / Pongs: %d\n", count)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte(status)); err != nil {
		log.Printf("writing response: %v", err)
	}
	fmt.Printf("%s -> %d\n", r.Pattern, count)
}
