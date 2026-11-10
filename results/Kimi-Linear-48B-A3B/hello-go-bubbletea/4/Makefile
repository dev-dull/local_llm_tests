.PHONY: run build clean install test

run:
	go run main.go

build:
	go build -o hello-bubbletea main.go

clean:
	go clean
	rm -f hello-bubbletea

install:
	go install

test:
	go test -v

# Cross-platform build
build-all:
	GOOS=linux GOARCH=amd64 go build -o hello-bubbletea-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o hello-bubbletea-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o hello-bubbletea-windows-amd64.exe main.go

# Development helpers
deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...