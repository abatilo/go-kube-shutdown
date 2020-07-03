package main

import (
	"log"
	"net/http"
	"time"

	shutdown "github.com/abatilo/go-kube-shutdown"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request, waiting...")
		time.Sleep(3 * time.Second)
		w.Write([]byte("Hello world"))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server...")
	shutdown.StartSafeServer(srv, &shutdown.LivenessOptions{
		Path: "/tmp/live",
	})
	log.Println("Server stopped...")
}
