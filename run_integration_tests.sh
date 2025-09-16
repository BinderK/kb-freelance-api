#!/bin/bash

# Integration Test Runner for kb-freelance-api
# This script runs integration tests with proper environment configuration

echo "🧪 Running Integration Tests..."
echo "================================"

# Check if required environment variables are set
if [ -z "$PYTHON_EXEC_PATH" ]; then
    echo "⚠️  PYTHON_EXEC_PATH not set. Using default: python3"
    echo "   To use a specific Python environment, set:"
    echo "   export PYTHON_EXEC_PATH='/path/to/your/python'"
    echo ""
fi

if [ -z "$TIME_TRACKER_PATH" ]; then
    echo "⚠️  TIME_TRACKER_PATH not set. Using default: ../kb-tt-cli"
    echo "   To use a specific path, set:"
    echo "   export TIME_TRACKER_PATH='/path/to/kb-tt-cli'"
    echo ""
fi

if [ -z "$INVOICE_GEN_PATH" ]; then
    echo "⚠️  INVOICE_GEN_PATH not set. Using default: ../kb-invoice-gen-cli"
    echo "   To use a specific path, set:"
    echo "   export INVOICE_GEN_PATH='/path/to/kb-invoice-gen-cli'"
    echo ""
fi

echo "Running integration tests with current configuration..."
echo ""

# Run integration tests
go test -v ./internal/services -run Integration
go test -v ./internal/api -run Integration

echo ""
echo "================================"
echo "✅ Integration tests completed!"
echo ""
echo "To run with custom paths:"
echo "  PYTHON_EXEC_PATH='/path/to/python' \\"
echo "  TIME_TRACKER_PATH='/path/to/kb-tt-cli' \\"
echo "  INVOICE_GEN_PATH='/path/to/kb-invoice-gen-cli' \\"
echo "  ./run_integration_tests.sh"



