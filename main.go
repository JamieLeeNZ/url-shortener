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

	http.HandleFunc("/shorten", s.ShortenHandler)

	http.HandleFunc("/delete/", s.DeleteHandler)

	http.HandleFunc("/", s.GetHandler)

	port := ":8080"
	log.Printf("Starting server at http://localhost%s/health\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
