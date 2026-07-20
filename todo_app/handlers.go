package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"text/template"
	"time"
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
				log.Printf("failed to refresh image: %v", err)
			}
		}()
	}

	tdh.mu.Unlock()

	http.ServeFile(w, r, imagePath)
}

func (tdh *todoAppHandler) handleRoot(w http.ResponseWriter, _ *http.Request) {
	resp, err := http.Get(tdh.todoBackendURL + "/todos")
	if err != nil {
		log.Printf("Failed send request: error: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected todo backend response status: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var req todoResponse
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Get todos: mailformed backend response body: %v", err)
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
		log.Printf("render page: %v", err)
	}
}

func (tdh *todoAppHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if title == "" || len(title) > 140 {
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
		http.Error(w, "backend unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated &&
		resp.StatusCode != http.StatusOK {
		http.Error(w, "backend error", http.StatusBadGateway)
		return
	}

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
