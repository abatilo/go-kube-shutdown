package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// StartSafeServer wraps the server to shutdown gracefully. This implementation handles SIGINT and SIGTERM
//
// Implementation heavily inspired by:
//
// https://medium.com/over-engineering/graceful-shutdown-with-go-http-servers-and-kubernetes-rolling-updates-6697e7db17cf
func StartSafeServer(server *http.Server, livenessMarkerPath string) error {
	livenessFile, err := os.Create(livenessMarkerPath)
	if err != nil {
		return fmt.Errorf("Failed to create liveness marker file: %v", err)
	}
	defer os.Remove(livenessFile.Name())

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// SIGINT signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// SIGTERM signal sent from Kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		server.Shutdown(context.Background())
		close(idleConnsClosed)
	}()

	err = server.ListenAndServe()
	<-idleConnsClosed
	return err
}
