.PHONY: build run clean install-deps demo

# Build the application
build:
	go build -o hello-go-bubbletea main.go

# Install dependencies
install-deps:
	go mod tidy
	go mod download

# Run the application
run: build
	./hello-go-bubbletea

# Clean build artifacts
clean:
	go clean
	rm -f hello-go-bubbletea

# Run with demo high score
demo:
	BUBBLE_HIGH_SCORE=999 go run main.go

# Development mode with live reload (requires air)
dev:
	air

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the application"
	@echo "  run      - Build and run the application"
	@echo "  clean    - Clean build artifacts"
	@echo "  install-deps - Install and tidy dependencies"
	@echo "  demo     - Run with demo high score"
	@echo "  dev      - Run in development mode (requires air)"
	@echo "  help     - Show this help message"