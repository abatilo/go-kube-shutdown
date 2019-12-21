package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/abatilo/go-kube-shutdown/pkg/ready"
	"github.com/abatilo/go-kube-shutdown/pkg/shutdown"
)

func main() {
	readyChecks := ready.NewChecks()
	readyChecks.Add("passes", func() error {
		return nil
	})
	readyChecks.Add("fails", func() error {
		return errors.New("Failure")
	})
	readyChecks.Add("google", ready.HTTPGet("https://www.google.com"))

	healthcheckServer := &http.Server{
		// Run on a different port that isn't exposed to the world
		Addr:    ":9091",
		Handler: readyChecks,
	}
	// Run alongside your main web server
	go healthcheckServer.ListenAndServe()

	// Use a default router for serving requests
	router := http.NewServeMux()

	// Simulate a long running request
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		log.Printf("Waiting 10 seconds for request %s to finish", id)
		<-time.After(time.Second * 10)
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
