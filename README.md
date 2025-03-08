# Wallet Service - Financial Asset Management Application Services

This application allows you to manage your different financial assets (bank accounts, stock market investments, cryptocurrencies, etc.) on a single platform.

## Features

-   **Cash Assets**: Bank accounts, cash
-   **Investment Assets**: Stocks, mutual funds, bonds, futures (VIOP - Turkish Derivatives Exchange)
-   **Crypto Assets**: Cryptocurrencies, NFTs, DeFi assets
-   **Other Assets**: Gold/silver, real estate, debts/receivables
-   **Asset Tracking**: View the total value of all your assets and their performance over time.
-   **Transaction Recording**: Keep a record of all your financial transactions.
-   **API Integration**: Automatic data synchronization with bank, stock exchange, and cryptocurrency exchange APIs.

## Technologies

-   **Backend**: Go (Golang)
-   **Web Framework**: Fiber
-   **Database**: PostgreSQL
-   **Authentication**: JWT

## Installation

### Prerequisites

-   Go 1.21+
-   PostgreSQL
-   git

### Database Setup

```bash
# Create a PostgreSQL database
createdb wallet

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

5.  Run the Worker Service (optional, for background tasks):

```bash
make run-worker
# or directly:
go run cmd/worker/main.go
```

## API Endpoints

### Authentication

-   `POST /api/auth/register` - Register a new user
-   `POST /api/auth/login` - User login
-   `GET /api/auth/me` - View current user information
-   `PUT /api/auth/me` - User update
-   `PUT /api/auth/change-password` - User change password
-   `DELETE /api/auth/me` - User delete

### Accounts

-   `GET /api/accounts` - List all accounts
-   `POST /api/accounts` - Create a new account
-   `GET /api/accounts/{id}` - View a specific account
-   `PUT /api/accounts/{id}` - Update account information
-   `DELETE /api/accounts/{id}` - Delete an account
-   `POST /api/accounts/{id}/credentials` - Set account API credentials
-   `GET /api/accounts/types` - List available account types

### Assets

-   `GET /api/assets` - List all assets
-   `POST /api/assets` - Add a new asset
-   `GET /api/assets/{id}` - View a specific asset
-   `PUT /api/assets/{id}` - Update asset information
-   `DELETE /api/assets/{id}` - Delete an asset
-   `GET /api/assets/types` - List available asset types

### Transactions

-   `GET /api/transactions` - List all transactions
-   `POST /api/transactions` - Add a new transaction
-   `GET /api/transactions/{id}` - View a specific transaction
-   `PUT /api/transactions/{id}` - Update transaction information
-   `DELETE /api/transactions/{id}` - Delete a transaction
-   `GET /api/transactions/types` - List available transaction types

## Project Structure

```
wallet-service/
├── cmd/                # Main executable files (main packages)
│   ├── api/            # API server (REST, gRPC, etc.)
│   │   └── main.go
│   └── worker/         # Workers that handle background tasks (e.g., periodic data synchronization) (Not fully planned yet. In the future...)
│       └── main.go
├── internal/           # Application-specific, non-exported (private) code
│   ├── app/            # Application layer (business logic services)
│   │   ├── users/       # User management service
│   │   ├── accounts/    # Account management service
│   │   ├── assets/      # Asset management service
│   │   ├── transactions/ # Transaction management service
│   │   └── ...
│   ├── pkg/            # Utility packages that can be used by different parts of the application
│   │   ├── auth/       # Authentication and authorization
│   │   ├── httpclient/  # HTTP client configuration and utility functions
│   │   ├── logger/     # Logging
│   │   ├── config/     # Configuration management
│   │   └── ...
│   ├── platform/       # Integrations with external services (3rd-party APIs, database)
│   │   ├── database/   # Database connection and operations (if using an ORM, model definitions might be here)
│   │   ├── bankapi/   # Integration with bank APIs
│   │   ├── exchangeapi/ # Integration with stock exchange APIs
│   │   ├── cryptoapi/  # Integration with cryptocurrency exchange APIs
│   │   └── ...
├── pkg/                # General-purpose (public) packages that can be used in other projects (optional)
│   ├── api/            # API definitions (Protobuf, OpenAPI, etc.)
│   └── ...
├── domain/             # Domain objects and rules (DDD)
│   ├── user/          # User model and related business rules
│   │   ├── user.go
│   │   ├── repository.go  # Interface for accessing user data
│   │   └── service.go     # Business logic related to users (optional, can be combined with `internal/app`)
│   ├── account/       # Account model and related business rules
│   ├── asset/         # Asset model and related business rules
│   ├── transaction/  # Transaction model and related business rules
│   └── ...
├── infrastructure/   # Infrastructure layer (DDD) - Interaction with database, external services, etc. (optional)
│   ├── persistence/ # Database operations (Repository implementations)
│   │   ├── userrepo/   # UserRepository implementation
│   │   ├── accountrepo/
│   │   └── ...
│   ├── external/      # Interaction with external services
│   │   ├── bank/      # Bank API client
│   │   ├── exchange/  # Stock exchange API client
│   │   └── ...
├── api/           # Presentation Layer - API handlers, request/response models
│  ├── handlers/  # HTTP handlers
│  ├── models/   # Request/response data structures
│  └── middleware/ # Middlewares (authentication, logging, etc.)
├── scripts/            # Utility scripts (database migrations, deployment, etc.)
├── deployments/        # Deployment configurations (Docker, Kubernetes, etc.)
├── configs/            # Configuration files (environment-specific settings)
├── test/               # Tests (unit, integration, e2e)
├── Makefile            # Shortcuts for common tasks (build, test, deploy, etc.)
└── go.mod, go.sum    # Go modules
```

## Development

### Creating a New Migration

```bash
make migrate-create name=migration_name
```

### Running Tests

```bash
make test
```

### Building the Application

```bash
make build
```