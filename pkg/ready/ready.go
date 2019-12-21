// Package ready makes it easy to verify that your applications is ready to
// receive network traffic
//
// Readiness checks are for Kubernetes to determine whether or not your application is ready to be routed network traffic. Common checks for readiness checks might be validating a connection with your database, or perhaps sending a request to a different service's readiness or ping endpoints.
//
// Read more about readiness checks here:
//
// https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
//
// Implementation heavily inspired by:
//
// https://github.com/heptiolabs/healthcheck
package ready

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Check is a health/readiness check.
type Check func() error

// Handler does things
type Handler interface {
	http.Handler

	// AddReadinessCheck adds a check that indicates that this instance of the
	// application is currently unable to serve requests because of an upstream
	// or some transient failure. If a readiness check fails, this instance
	// should no longer receiver requests, but should not be restarted or
	// destroyed.
	Add(name string, check Check)
}

// basicHandler is a basic Handler implementation.
type basicHandler struct {
	http.ServeMux
	checksMutex     sync.RWMutex
	readinessChecks map[string]Check
}

// NewChecks creates a new handler where you can register healthchecks
func NewChecks() Handler {
	h := &basicHandler{
		readinessChecks: make(map[string]Check),
	}
	h.HandleFunc("/ready", h.handle)
	return h
}

func (s *basicHandler) Add(name string, check Check) {
	s.checksMutex.Lock()
	defer s.checksMutex.Unlock()
	s.readinessChecks[name] = check
}

func (s *basicHandler) collectChecks(resultsOut map[string]string, statusOut *int) {
	s.checksMutex.RLock()
	defer s.checksMutex.RUnlock()
	for name, check := range s.readinessChecks {
		if err := check(); err != nil {
			*statusOut = http.StatusServiceUnavailable
			resultsOut[name] = err.Error()
		} else {
			resultsOut[name] = "OK"
		}
	}
}

func (s *basicHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	checkResults := make(map[string]string)
	status := http.StatusOK
	s.collectChecks(checkResults, &status)

	// write out the response code and content type header
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// v for verbose. unless ?v=1, return an empty body. Kubernetes only cares
	// about the HTTP status code, so we won't waste bytes on the full body.
	if r.URL.Query().Get("v") != "1" {
		w.Write([]byte("{}\n"))
		return
	}

	bytes, err := json.Marshal(checkResults)
	if err == nil {
		w.Write(bytes)
	}
}
