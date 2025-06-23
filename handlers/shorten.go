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

func (s *Server) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req models.URLShortenRequest
	db := s.store

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Original URL:", req.Original)

	if req.Original == "" {
		http.Error(w, "Original URL is required", http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(req.Original)
	if err != nil {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	key := generateRandomKey(6)
	for db.Contains(key) {
		key = generateRandomKey(6)
	}

	db.Set(key, req.Original)

	resp := models.URLShortenResponse{Key: key}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI required.", http.StatusBadRequest)
		return
	}

	original, ok := s.store.Get(key)
	if !ok {
		http.Error(w, "Invalid URL.", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
