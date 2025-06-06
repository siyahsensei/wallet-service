#!/bin/bash

echo "🚀 Starting Wallet Service with Swagger Documentation..."
echo ""
echo "📚 Swagger UI will be available at: http://localhost:8080/swagger/"
echo "🏥 Health check endpoint: http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Start the application
go run ./cmd/api/main.go 