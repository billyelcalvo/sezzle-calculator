package main

import (
	"log"
	"net/http"
	"time"

	"github.com/diegomora/sezzle-calculator/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", handlers.Calculate)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("calculator API listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
