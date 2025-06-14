install:
	go build -o $(GOPATH)/bin/gf .

test:
	go test -v ./...