package main

import (
	"log"
	"os"

	"bootstrap-server/internal/bootstrap"
)

func main() {
	// Konfigurasi server (misalnya port diambil dari ENV atau default ke 4000)
	port := os.Getenv("BOOTSTRAP_PORT")
	if port == "" {
		port = "4000"
	}

	server, err := bootstrap.NewServer(port, "data/peers.db")
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	log.Printf("Bootstrap server listening on port %s\n", port)

	// Jalankan server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
