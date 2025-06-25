package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JamieLeeNZ/url-shortener/handlers"
	"github.com/JamieLeeNZ/url-shortener/store"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	postgresStore, err := store.NewPostgresStore(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer postgresStore.Close()

	memStore := store.NewMemoryStore()

	s := handlers.NewServer(memStore)

	http.HandleFunc("/health", s.HealthHandler)

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
