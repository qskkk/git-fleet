# Get version from git tag or commit
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/qskkk/git-fleet/config.Version=$(VERSION)"

# Build binary with version
build:
	go build $(LDFLAGS) -o gf .

install:
	go build $(LDFLAGS) -o $(GOPATH)/bin/gf .

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/gf-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/gf-windows-amd64.exe .

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -f gf
	rm -rf dist/

.PHONY: build install build-all test test-cover clean