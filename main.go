package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "OK")
    })

    port := ":8080"
    log.Printf("Starting server at http://localhost%s/health\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}