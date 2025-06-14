install:
	go build -o $(GOPATH)/bin/gf .

test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out