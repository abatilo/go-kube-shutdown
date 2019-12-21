# go-kube-shutdown
An opinionated library for handling Kubernetes readiness, liveness, and
shutdown concepts as a first class citizen.

## Examples
TODO

## Motivation
I've seen a lot of libraries that get parts of the application shutdown
incorrect, or the difference in how you're supposed to handle readiness vs
liveness checks in a Golang web service application.

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
