package main

import (
	"log"
	"net/http"

	"github.com/JamieLeeNZ/url-shortener/handlers"
	"github.com/JamieLeeNZ/url-shortener/store"
)

func main() {
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
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	port := ":8080"
	log.Printf("Starting server at http://localhost%s/health\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
