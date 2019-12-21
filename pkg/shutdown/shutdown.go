package shutdown

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// StartSafeServer wraps the server to shutdown gracefully. This implementation
// handles SIGINT and SIGTERM cleanly. Pass in your *http.Server which will be
// started for you. Second, pass in a path for where your liveness marker file
// should live. This can be used by your Kubernetes liveness probe to see
// whether the server has started. We will create a file at that location right
// before the server starts and we will delete the file right after the server
// shutsdown. This lets you decouple the requirement for http traffic to be
// routable to the service in order to figure out if the application process
// has been stopped.
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
