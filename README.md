# go-kube-shutdown
[![GoDoc](https://godoc.org/github.com/abatilo/go-kube-shutdown?status.svg)](http://godoc.org/github.com/abatilo/go-kube-shutdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/abatilo/go-kube-shutdown)](https://goreportcard.com/report/github.com/abatilo/go-kube-shutdown)

An opinionated library for handling Kubernetes readiness, liveness, and
shutdown concepts as a first class citizen.

## Usage
Installation:
```
go get -u github.com/abatilo/go-kube-shutdown
```

## Examples
Full examples can be found in the [examples](./examples/) directory.

### Readiness / Healthchecks
Import:
```
import "github.com/abatilo/go-kube-shutdown/pkg/ready"
```

Add a server that runs on a different port to respond to readiness checks
```go
readyChecks := ready.NewChecker()
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
```

### Graceful shutdown
Import:
```
import "github.com/abatilo/go-kube-shutdown/pkg/shutdown"
```

Create your server ahead of time and start it by running `shutdown.StartSafeServer`.
```go
server := &http.Server{
	Addr:    ":9090",
}

log.Printf("Starting a server that will shutdown safely")
livenessFileMarker := "/tmp/liveness"
err := shutdown.StartSafeServer(server, livenessFileMarker)
if err != http.ErrServerClosed {
	log.Printf("Server did not shutdown cleanly: %v", err)
}
log.Printf("Connections have drained from the server and the server has shutdown")
```

Your server will be started using `server.ListenAndServe()` but is now wrapped with a goroutine and channel that will ensure that on shutdown, we wait for connections to drain before actually exiting.

Check the [GoDoc](https://godoc.org/github.com/abatilo/go-kube-shutdown?status.svg) for more details.

## Motivation
I've seen a lot of libraries that misunderstand parts of the application
shutdown incorrectly, or misunderstand the difference in how you're supposed to
handle readiness vs liveness checks in a Golang web service application.

I've copied and pasted the same snippets of code for multiple projects now, and
I just want to wrap all of the things that I do into a single library, so that
I can bootstrap a new Golang microservice very quickly.

## Thanks
Takes inspiration and ideas from:
* https://github.com/InVisionApp/go-health
* https://github.com/heptiolabs/healthcheck
* https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
* https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
* https://medium.com/over-engineering/graceful-shutdown-with-go-http-servers-and-kubernetes-rolling-updates-6697e7db17cf
