# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Dummy Loader is a small Go HTTP server that exposes an API for generating CPU and memory load, used to experiment with resource management config in Docker/Kubernetes. Module: `github.com/sha1n/dummy-loader`, Go 1.26.

The buildable package is one level down from the module root: `github.com/sha1n/dummy-loader/server` (package `main`, entrypoint `server/bootstrap.go`). There's no `cmd/`/`pkg/`/`internal/` layout — packages live directly under `server/` (`http/`, `loaders/`, `sys/`, `utils/`, `web/`).

## Build, test, lint

All via the root `Makefile`:
- `make test` — runs `go test -v` across all packages
- `make format` — `gofmt -s -w server`
- `make lint` — builds a pinned `golangci-lint` binary from the `tools` submodule into `.bin/`, then runs it (config: `.golangci.yml`)
- `make build-bin` — builds the Linux binary to `server/bin/dummy-loader`
- `make build-docker` — builds the Docker image; **requires `build-bin` to have already run**, since the Dockerfile just copies a prebuilt binary rather than doing a multi-stage Go build. Use `make build` to run both in order.
- `make setup` — points git at `.git-hooks/` (`core.hooksPath`), enabling a pre-commit hook that blocks commits with unformatted Go files. Not enabled by default on a fresh clone.

Tests are colocated with source (standard Go convention) and use `testify/assert`. `server/http/server_test.go` spins up a real HTTP server on a random port and issues real requests rather than mocking.

## Notes

- `.travis.yml` targets Go 1.15 and is stale relative to go.mod's Go 1.25 — don't treat it as the source of truth for CI behavior.
- `Gopkg.toml`/`Gopkg.lock` are leftovers from the pre-modules `dep` tool and are unused by the current build.
- `make lint` runs clean as of this writing (the pre-existing `errcheck`/`revive` backlog noted earlier has been resolved).
