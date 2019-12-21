package ready_test

import (
	"net/http"

	"github.com/abatilo/go-kube-shutdown/pkg/ready"
)

func Example() {
	readyChecks := ready.NewChecks()
	readyChecks.Add("google", ready.HTTPGet("https://www.google.com"))

	healthcheckServer := &http.Server{
		Addr:    ":9091",
		Handler: readyChecks,
	}
	go healthcheckServer.ListenAndServe()
}

func Example_customCheck() {
	readyChecks := ready.NewChecks()
	readyChecks.Add("passes", func() error {
		return nil
	})

	healthcheckServer := &http.Server{
		Addr:    ":9091",
		Handler: readyChecks,
	}
	go healthcheckServer.ListenAndServe()
}
