# Wallet Service - Financial Asset Management Application Services

This application allows you to manage your different financial assets (bank accounts, stock market investments, cryptocurrencies, etc.) on a single platform.

## Technologies

-   **Backend**: Go (Golang)
-   **Web Framework**: Fiber
-   **Database**: PostgreSQL
-   **Authentication**: JWT
-   **API Documentation**: Swagger/OpenAPI

## API Documentation

This project includes automatically generated Swagger documentation for all API endpoints.

### Accessing Swagger UI

Once the application is running, you can access the interactive Swagger UI at:

```
http://localhost:8080/swagger/
```

### Regenerating Documentation

To regenerate the Swagger documentation after making changes to the API:

```bash
# Install swag CLI tool (one time setup)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate/update documentation
make swagger
```

### Authentication in Swagger

For endpoints that require authentication:
1. First, use the `/api/auth/login` or `/api/auth/register` endpoint to get a JWT token
2. Click the "Authorize" button in Swagger UI
3. Enter `Bearer <your-jwt-token>` in the Authorization field
4. Now you can test authenticated endpoints

## Installation

### Prerequisites

-   Go 1.23+
-   PostgreSQL
-   git

### Database Setup

```bash
# Create a PostgreSQL database
createdb wallety

# Install golang-migrate for running migrations
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
make migrate-up
```

### Running the Application

1.  Clone the repository:

```bash
git clone https://github.com/siyahsensei/wallet-service.git
cd wallet-service
```

2.  Install dependencies:

```bash
go mod download
```

3.  Create a `.env` file for configuration:

```bash
cp .env.example .env
# Edit the .env file and configure the necessary settings.
```

4.  Run the API Server:

```bash
make run-api
# or directly:
go run cmd/api/main.go
```