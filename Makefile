# Get version from git tag or commit
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X 'github.com/qskkk/git-fleet/internal/pkg/version.Version=$(VERSION)' -X 'github.com/qskkk/git-fleet/internal/pkg/version.BuildDate=$(BUILD_DATE)' -X 'github.com/qskkk/git-fleet/internal/pkg/version.GitCommit=$(GIT_COMMIT)'"

# Build binary with version
build:
	go build $(LDFLAGS) -o gf ./cmd/gf

install:
	go build $(LDFLAGS) -o $(GOPATH)/bin/gf ./cmd/gf

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-darwin-amd64 ./cmd/gf
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/gf-darwin-arm64 ./cmd/gf
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-linux-amd64 ./cmd/gf
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-windows-amd64.exe ./cmd/gf

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -f gf
	rm -rf dist/

.PHONY: build install build-all test test-cover clean