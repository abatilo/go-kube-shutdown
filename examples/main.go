package main

import (
	"log"
	"net/http"
	"time"

	"github.com/abatilo/go-kube-shutdown/pkg/shutdown"
)

func main() {
	// Use a default router for serving requests
	router := http.NewServeMux()

	// Simulate a long running request
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		<-time.After(time.Second * 5)
		w.Write([]byte("pong"))
	})

	server := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	log.Printf("Starting a server that will shutdown safely")
	livenessFileMarker := "/tmp/liveness"
	err := shutdown.StartSafeServer(server, livenessFileMarker)
	if err != http.ErrServerClosed {
		log.Printf("Server did not shutdown cleanly: %v", err)
	}
	log.Printf("Connections have drained from the server and the server has shutdown")
}
