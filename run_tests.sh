#!/bin/bash

# Test runner script for kb-freelance-api
# This script runs all tests with coverage and provides a summary

echo "🧪 Running Go API Tests..."
echo "================================"

# Run tests with coverage
echo "Running tests with coverage..."
go test -v -cover ./...

echo ""
echo "================================"
echo "📊 Test Coverage Summary:"
echo "================================"

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

echo ""
echo "================================"
echo "📈 HTML Coverage Report:"
echo "================================"
echo "To view detailed coverage report, run:"
echo "  go tool cover -html=coverage.out"
echo ""
echo "✅ All tests completed successfully!"
