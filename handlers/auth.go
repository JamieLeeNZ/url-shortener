package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JamieLeeNZ/url-shortener/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func InitOAuth() {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURL := "http://localhost:8080/auth/google/callback"

	if clientID == "" || clientSecret == "" {
		log.Fatal("Missing CLIENT_ID or CLIENT_SECRET")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func generateOauthState(w http.ResponseWriter) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Println("failed to generate random state:", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return ""
	}

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Expires:  time.Now().Add(20 * time.Minute),
		HttpOnly: true,
		Secure:   false, // Change to true in production with HTTPS
	}
	http.SetCookie(w, &cookie)

	return state
}

func (s *Server) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateOauthState(w)
	if state == "" {
		return
	}
	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "State cookie missing", http.StatusBadRequest)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println("getUserDataFromGoogle error:", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var gUser GoogleUser
	if err := json.Unmarshal(data, &gUser); err != nil {
		log.Println("json unmarshal error:", err)
		http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:      gUser.ID,
		Email:   gUser.Email,
		Name:    gUser.Name,
		Picture: gUser.Picture,
	}

	savedUser, err := s.userStore.GetOrCreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	if err := s.createSession(w, r.Context(), savedUser); err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedUser)
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
