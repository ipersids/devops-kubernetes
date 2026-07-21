package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type todoResponse struct {
	Data []todo `json:"data"`
}

type pageData struct {
	Title      string
	Todos      []todo
	BackendURL string
}

type todoAppHandler struct {
	volumeDir      string
	imageName      string
	tmpl           *template.Template
	mu             sync.Mutex
	refreshing     bool
	todoBackendURL string
}

func (tdh *todoAppHandler) handleImage(w http.ResponseWriter, r *http.Request) {
	imagePath := tdh.volumeDir + tdh.imageName
	tdh.mu.Lock()

	info, err := os.Stat(imagePath)
	if err != nil {
		tdh.mu.Unlock()
		logger(r, http.StatusInternalServerError).Msg("image not found")
		http.Error(w, "image not found", http.StatusInternalServerError)
		return
	}

	expired := time.Since(info.ModTime()) >= 10*time.Minute

	if expired && !tdh.refreshing {
		tdh.refreshing = true

		go func() {
			defer func() {
				tdh.mu.Lock()
				tdh.refreshing = false
				tdh.mu.Unlock()
			}()

			if err := downloadAndReplaceImage(tdh.volumeDir, tdh.imageName); err != nil {
				logger(r, http.StatusInternalServerError).Err(err).Msg("failed to refresh image")
			}
		}()
	}

	tdh.mu.Unlock()

	logger(r, http.StatusOK).Msg("successful request")
	http.ServeFile(w, r, imagePath)
}

func (tdh *todoAppHandler) handleRoot(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(tdh.todoBackendURL + "/todos")
	if err != nil {
		logger(r, http.StatusInternalServerError).Err(err).Msg("failed to send request")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger(r, http.StatusInternalServerError).Err(err).Msg("unexpected todo backend response status")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var req todoResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		logger(r, http.StatusInternalServerError).Err(err).Msg("mailformed backend response body")
		http.Error(w, "Failed fetch todos", http.StatusInternalServerError)
		return
	}

	todoList := req.Data
	data := pageData{
		Title:      "My day",
		Todos:      todoList,
		BackendURL: tdh.todoBackendURL,
	}

	if err := tdh.tmpl.Execute(w, data); err != nil {
		logger(r, http.StatusInternalServerError).Err(err).Msg("failed to execute template")
		return
	}

	logger(r, http.StatusOK).Msg("successful request")
}

func (tdh *todoAppHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logger(r, http.StatusBadRequest).Err(err).Msg("invalid form")
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if title == "" || len(title) > 140 {
		logger(r, http.StatusBadRequest).Msg("title should be between 1 and 140 bytes")
		http.Error(w, "Title should be between 1 and 140 bytes", http.StatusBadRequest)
		return
	}

	resp, err := http.PostForm(
		tdh.todoBackendURL+"/todos",
		url.Values{
			"title": {title},
		},
	)
	if err != nil {
		logger(r, http.StatusBadGateway).Msg("backend unavailable")
		http.Error(w, "backend unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK {
		logger(r, http.StatusBadGateway).Msg("backend error")
		http.Error(w, "backend error", http.StatusBadGateway)
		return
	}

	logger(r, http.StatusSeeOther).Msg("task created, redirecting to main page")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HELPERS

func downloadAndReplaceImage(path string, fileName string) error {
	tmp := path + ".tmp"

	resp, err := http.Get("https://picsum.photos/1200")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(tmp)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return os.Rename(tmp, path+fileName)
}

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
