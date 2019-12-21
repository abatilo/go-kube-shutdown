SHELL := /bin/bash

GO_VERSION = 1.13.4

.PHONY: help
help: ## View help information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Builds the project
	docker run --rm -it -v`pwd`:/src -w /src \
		golang:$(GO_VERSION) go build ./...
