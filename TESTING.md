Testing esctl

This project uses Go’s standard testing tools. Unit tests exercise the HTTP layer with a local mock server and do not require a live Elasticsearch cluster.

Prerequisites

- Go 1.23.x installed and on your PATH
- A working network connection for the first `go mod download` (to fetch dependencies)

Getting Started

- Clone the repo and change into the project directory
- Ensure the correct Go version: `go version` should report 1.23.x
- Download modules (first run only): `go mod download`

Running Tests

- All packages: `go test ./...`
- Verbose output: `go test -v ./...`
- Single package: `go test ./es/index`
- Single test: `go test ./es/index -run TestUpdateIndexSettings -v`
- With the race detector: `go test -race ./...`

Coverage

- Package coverage: `go test -cover ./...`
- HTML report: `go test -coverprofile=cover.out ./... && go tool cover -html=cover.out`

Static Analysis

- Vet common issues: `go vet ./...`

Mocking HTTP

- Tests use `internal/testutil/mock.go` to spin up an in-memory `httptest.Server` and return a client bound to it. Example pattern:

  - Create mock server/client: `srv, cli := testutil.NewMockServer(json, "/path")`
  - Inject into runtime: `shared.SetClient(cli)`
  - Defer `srv.Close()` to clean up

Debug Logging

- Many HTTP helpers honor a global debug flag. When running commands via the CLI you can add `--debug` to print request/response status lines to stderr. This is not required for tests.

CI

- GitHub Actions run tests on Go 1.23.x. If you hit version‑parsing errors locally, verify your Go toolchain is 1.23.x and re‑run `go mod download`.

Common Issues

- Module download errors: ensure you have network access for the initial `go mod download`. Subsequent runs use the module cache.
- Toolchain mismatch: if you see messages about an invalid `go` directive, update your Go toolchain to 1.23.x.

