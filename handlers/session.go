package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/JamieLeeNZ/url-shortener/models"
	"github.com/google/uuid"
)

type contextKey string

const (
	sessionCookieName = "session_id"
	sessionPrefix     = "session:"
	sessionDuration   = 24 * time.Hour
	userContextKey    = contextKey("user")
)

func (s *Server) createSession(w http.ResponseWriter, ctx context.Context, user models.User) error {
	sessionID := uuid.New().String()

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.redisClient.Set(ctx, sessionPrefix+sessionID, data, sessionDuration).Err()
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Expires:  time.Now().Add(sessionDuration),
		HttpOnly: true,
		Secure:   false, // set to true in production with HTTPS
		Path:     "/",
	})

	return nil
}

func (s *Server) getSessionUser(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil, err
	}

	data, err := s.redisClient.Get(r.Context(), sessionPrefix+cookie.Value).Result()
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	s.redisClient.Expire(r.Context(), sessionPrefix+cookie.Value, sessionDuration)

	return &user, nil
}

func (s *Server) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.getSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func GetCurrentUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil {
		s.redisClient.Del(r.Context(), sessionPrefix+cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
