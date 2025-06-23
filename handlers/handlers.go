package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/JamieLeeNZ/url-shortener/models"
	"github.com/JamieLeeNZ/url-shortener/store"
)

type Server struct {
	store *store.MemoryStore
}

func NewServer(db *store.MemoryStore) *Server {
	return &Server{
		store: db,
	}
}

func generateRandomKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "OK"}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This is a POST method only.", http.StatusMethodNotAllowed)
		return
	}

	db := s.store

	req, err := parseAndValidateURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var key string

	if k, found := db.GetKeyFromOriginal(req.Original); found {
		key = k
	} else {
		key = generateRandomKey(6)
		for db.ContainsKey(key) {
			key = generateRandomKey(6)
		}
		db.Set(key, req.Original)
	}

	resp := models.URLShortenResponse{Key: key}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "This is a GET method only.", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI required.", http.StatusBadRequest)
		return
	}

	original, ok := s.store.GetOriginalFromKey(key)
	if !ok {
		http.Error(w, "Invalid URL.", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "This is a PUT method only.", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI required.", http.StatusBadRequest)
		return
	}

	req, err := parseAndValidateURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	success := s.store.Update(key, req.Original)
	if !success {
		http.Error(w, "Key not found or new URL already mapped to a different key", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "This is a DELETE method only.", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI required.", http.StatusBadRequest)
		return
	}

	if !s.store.ContainsKey(key) {
		http.Error(w, "Invalid URL.", http.StatusNotFound)
		return
	}

	s.store.Delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func parseAndValidateURL(r *http.Request) (models.URLShortenRequest, error) {
	var req models.URLShortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, fmt.Errorf("invalid JSON")
	}

	if req.Original == "" {
		return req, fmt.Errorf("original URL is required")
	}

	if _, err := url.ParseRequestURI(req.Original); err != nil {
		return req, fmt.Errorf("invalid URL format")
	}

	return req, nil
}
