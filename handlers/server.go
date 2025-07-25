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
	"github.com/redis/go-redis/v9"
)

type Server struct {
	urlStore    store.URLStore
	userStore   store.UserStore
	redisClient *redis.Client
}

func NewServer(urlStore store.URLStore, userStore store.UserStore, redisClient *redis.Client) *Server {
	return &Server{
		urlStore:    urlStore,
		userStore:   userStore,
		redisClient: redisClient,
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
		http.Error(w, "this is a POST method only", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	user := GetCurrentUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	db := s.urlStore

	req, err := parseAndValidateURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var key string

	if k, _, found := db.GetKeyFromOriginal(ctx, req.Original); found {
		key = k
	} else {
		key = generateRandomKey(6)
		for db.ContainsKey(ctx, key) {
			key = generateRandomKey(6)
		}

		if err := db.Set(ctx, key, req.Original, user.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	resp := models.URLShortenResponse{Key: key}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this is a GET method only", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI key is required", http.StatusBadRequest)
		return
	}

	original, _, ok := s.urlStore.GetOriginalFromKey(ctx, key)
	if !ok {
		http.Error(w, "invalid URL", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "this is a PUT method only", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	user := GetCurrentUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI key is required", http.StatusBadRequest)
		return
	}

	req, err := parseAndValidateURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, existingUserID, found := s.urlStore.GetKeyFromOriginal(ctx, req.Original)
	if found && existingUserID != user.ID {
		http.Error(w, "forbidden: you do not own this URL", http.StatusForbidden)
		return
	}

	success := s.urlStore.Update(ctx, key, req.Original)
	if !success {
		http.Error(w, "key not found or new URL already mapped to a different key", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (s *Server) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "this is a DELETE method only", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	user := GetCurrentUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	key := strings.TrimPrefix(r.URL.Path, "/")
	if key == "" {
		http.Error(w, "URI key is required", http.StatusBadRequest)
		return
	}

	_, existingUserID, found := s.urlStore.GetOriginalFromKey(ctx, key)
	if !found {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if existingUserID != user.ID {
		http.Error(w, "forbidden: you do not own this URL", http.StatusForbidden)
		return
	}

	if !s.urlStore.Delete(ctx, key) {
		http.Error(w, "failed to delete URL", http.StatusInternalServerError)
		return
	}
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

func (s *Server) ListUserLinks(w http.ResponseWriter, r *http.Request) {
	user := GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	urls, err := s.userStore.GetURLsByUserID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(urls)
}
