.PHONY: build build-tui build-menubar clean install test lint run run-tui run-menubar

# Default target
all: build

# Build both binaries
build: build-tui build-menubar

# Build TUI version
build-tui:
	go build -o bin/jean ./cmd/tui

# Build menu bar version
build-menubar:
	go build -o bin/jean-menubar ./cmd/menubar

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f jean jean-menubar

# Install to /usr/local/bin
install: build
	cp bin/jean-menubar /usr/local/bin/jean-menubar
	cp bin/jean /usr/local/bin/jean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Run menu bar app
run: run-menubar

run-menubar: build-menubar
	./bin/jean-menubar

# Run TUI app
run-tui: build-tui
	./bin/jean

# Development: rebuild and restart menu bar
dev:
	pkill -f jean-menubar || true
	$(MAKE) build-menubar
	./bin/jean-menubar &

# Tidy dependencies
tidy:
	go mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build both TUI and menu bar binaries"
	@echo "  build-tui      - Build TUI version only"
	@echo "  build-menubar  - Build menu bar version only"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install binaries to /usr/local/bin"
	@echo "  test           - Run tests"
	@echo "  lint           - Run linter"
	@echo "  run            - Build and run menu bar app"
	@echo "  run-tui        - Build and run TUI app"
	@echo "  dev            - Restart menu bar app (for development)"
	@echo "  tidy           - Tidy go modules"
