package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JamieLeeNZ/url-shortener/handlers"
	"github.com/JamieLeeNZ/url-shortener/store"
)

func main() {
	memStore := store.NewMemoryStore()

	s := handlers.NewServer(memStore)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	http.HandleFunc("/shorten", s.ShortenHandler)

	http.HandleFunc("/", s.GetHandler)

	port := ":8080"
	log.Printf("Starting server at http://localhost%s/health\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
