.PHONY: dev build run templ install clean

# Development with live reload
dev:
	~/go/bin/air

# Build commands
build: templ
	go build -o ./bin/server ./cmd/main.go

templ:
	~/go/bin/templ generate

# Run production build
run: build
	./bin/server

# Install dependencies
install:
	go mod download
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/air-verse/air@latest

# Clean build artifacts
clean:
	rm -rf ./bin ./tmp
