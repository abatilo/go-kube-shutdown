# go-kube-shutdown
[![GoDoc](https://godoc.org/github.com/abatilo/go-kube-shutdown?status.svg)](http://godoc.org/github.com/abatilo/go-kube-shutdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/abatilo/go-kube-shutdown)](https://goreportcard.com/report/github.com/abatilo/go-kube-shutdown)
[![license](https://img.shields.io/github/license/abatilo/go-kube-shutdown.svg)](https://github.com/abatilo/go-kube-shutdown/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/abatilo/go-kube-shutdown.svg)](https://github.com/abatilo/go-kube-shutdown/releases/latest)
[![GitHub release date](https://img.shields.io/github/release-date/abatilo/go-kube-shutdown.svg)](https://github.com/abatilo/go-kube-shutdown/releases)
![Build](https://github.com/abatilo/go-kube-shutdown/workflows/.github/workflows/build.yml/badge.svg)

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
```go
import "github.com/abatilo/go-kube-shutdown"
```

Add a server that runs on a different port to respond to readiness checks
```go
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
```

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
