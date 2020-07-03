// Package shutdown makes it easy to add graceful shutdown to your application
package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPath = "/live"
)

type httpServer interface {
	Shutdown(ctx context.Context) error
	ListenAndServe() error
}

// LivenessOptions lets you configure the safe http server
type LivenessOptions struct {
	// Path is where the file should be written to
	Path string
}

// StartSafeServer wraps the server to shutdown gracefully. This implementation
// handles SIGINT and SIGTERM cleanly. Pass in your *http.Server which will be
// started for you. Second, pass in configuration for where your liveness
// marker file should live. This can be used by your Kubernetes liveness probe
// to see whether the server has started. We will create a file at that
// location right before the server starts and we will delete the file right
// after the server shutsdown. This lets you decouple the requirement for http
// traffic to be routable to the service in order to figure out if the
// application process has been stopped.
//
// Implementation heavily inspired by:
//
// https://medium.com/over-engineering/graceful-shutdown-with-go-http-servers-and-kubernetes-rolling-updates-6697e7db17cf
func StartSafeServer(server httpServer, options *LivenessOptions) error {
	path := options.Path
	if path == "" {
		path = defaultPath
	}

	livenessFile, err := os.Create(path)

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
