# Wallet Service - Financial Asset Management Application Services

This application allows you to manage your different financial assets (bank accounts, stock market investments, cryptocurrencies, etc.) on a single platform.

## Technologies

-   **Backend**: Go (Golang)
-   **Web Framework**: Fiber
-   **Database**: PostgreSQL
-   **Authentication**: JWT

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