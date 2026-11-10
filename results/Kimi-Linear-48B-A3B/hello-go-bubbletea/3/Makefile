.PHONY: build run clean install-deps

# Build the application
build:
	go build -o hello-go-bubbletea main.go

# Run the application
run:
	go run main.go

# Install dependencies
install-deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -f hello-go-bubbletea

# Run with verbose output
run-verbose:
	go run -v main.go

# Cross-compile for different platforms
cross-compile:
	GOOS=linux GOARCH=amd64 go build -o hello-go-bubbletea-linux main.go
	GOOS=darwin GOARCH=amd64 go build -o hello-go-bubbletea-macos main.go
	GOOS=windows GOARCH=amd64 go build -o hello-go-bubbletea.exe main.go

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  install-deps  - Install and verify dependencies"
	@echo "  clean         - Remove build artifacts"
	@echo "  run-verbose   - Run with verbose output"
	@echo "  cross-compile - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  help          - Show this help message"