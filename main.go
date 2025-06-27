package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JamieLeeNZ/url-shortener/handlers"
	"github.com/JamieLeeNZ/url-shortener/store"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		log.Fatal("REDIS_ADDRESS is not set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		log.Fatal("REDIS_PASSWORD is not set")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal("GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET not set")
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	oauthStateString := os.Getenv("OAUTH_STATE_STRING")
	if oauthStateString == "" {
		log.Fatal("OAUTH_STATE_STRING not set")
	}

	postgresStore, err := store.NewPostgresStore(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgresStore.Close()

	redisStore, err := store.NewRedisStore(redisAddress, redisPassword, 0, 24*time.Hour)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisStore.Close()

	cachedStore, err := store.NewCachedStore(redisStore, postgresStore)
	if err != nil {
		log.Fatalf("Failed to create cached store: %v", err)
	}

	s := handlers.NewServer(cachedStore)

	http.HandleFunc("/health", s.HealthHandler)

	// http.HandleFunc("/login", handleGoogleLogin)
	// http.HandleFunc("/auth/google/callback", handleGoogleCallback)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.CreateHandler(w, r)
		case http.MethodGet:
			s.GetHandler(w, r)
		case http.MethodPut:
			s.UpdateHandler(w, r)
		case http.MethodDelete:
			s.DeleteHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	log.Printf("Starting server at http://localhost%s/health\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
