.PHONY: build test run migrate-up migrate-down

# Build the application
build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker

# Run tests
test:
	go test -v ./...

# Run the API server
run-api:
	go run ./cmd/api/main.go

# Run the worker
run-worker:
	go run ./cmd/worker/main.go

# Create DB migrations (requires golang-migrate)
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Run migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/wallet?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/wallet?sslmode=disable" down 