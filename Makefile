.PHONY: dev build run templ install clean db-up db-down db-migrate db-seed

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

# Database commands
db-up:
	docker compose up -d

db-down:
	docker compose down

db-migrate:
	docker exec -i showcase-go-db-1 psql -U postgres -d showcase < db/schema.sql

db-seed:
	docker exec -i showcase-go-db-1 psql -U postgres -d showcase < db/seed.sql

db-reset: db-migrate db-seed
