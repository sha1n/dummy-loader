NOCOL=\x1b[0m
GREEN=\x1b[32;01m
RED=\x1b[31;01m
YELLOW=\x1b[33;01m

define print_title
	@echo "---"
	@echo "--- $(GREEN)$1$(NOCOL)"
	@echo "---"
endef


default: get lint format test build-bin

get:
	go get -t ./...

test:
	$(call print_title, Running tests...)
	go test -v `go list ./...`


build: build-bin build-docker

build-bin:
	$(call print_title,Building binaries...)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server/bin/dummy-loader github.com/sha1n/dummy-loader/server

build-docker:
	$(call print_title,Building docker image...)
	docker build -t sha1n/dummy-loader server


prepare:
	$(call print_title,Preparing go dependencies...)
	dep ensure -v


format:
	$(call print_title,Formatting go sources...)
	gofmt -s -w server


lint:
	$(call print_title,Lint...)
	gofmt -d server


push-docker:
	$(call print_title,Publishing docker image...)
	docker push sha1n/dummy-loader:latest


run-docker:
	docker run -d -p 8080:8080 sha1n/dummy-loader


release: prepare format lint test build push-docker

setup:
	git config core.hooksPath .git-hooks