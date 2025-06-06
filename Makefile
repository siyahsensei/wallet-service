.PHONY: build test run migrate-up migrate-down swagger migrate-new-structure

# Build the application
build:
	go build -o bin/api ./cmd/api

# Run tests
test:
	go test -v ./...

# Run the API server
run-api:
	go run ./cmd/api/main.go

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Create DB migrations (requires golang-migrate)
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Run migrations up
migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/wallety?sslmode=disable" up

# Run migrations down
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/wallety?sslmode=disable" down

# Apply the new account structure migration
migrate-new-structure:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/wallety?sslmode=disable" up 1 