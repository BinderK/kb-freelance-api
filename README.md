# KB Freelance API

A Go-based REST API that integrates with the existing Python CLI applications for time tracking and invoice generation.

## Features

- **Time Tracking**: Start/stop timers, get status, view today's summary
- **Invoice Generation**: Generate PDF invoices with line items
- **REST API**: Clean RESTful interface for frontend integration
- **CORS Support**: Ready for React frontend integration
- **Environment Configuration**: Fully configurable via environment variables
- **Comprehensive Testing**: Unit tests, integration tests, and test coverage
- **Python CLI Integration**: Seamless integration with existing Python tools

## Architecture

```
React Frontend (Port 3000)
    ↓ HTTP/REST
Go API Server (Port 8080)
    ↓ Process execution
Python CLI Applications
    ↓ SQLAlchemy
SQLite Database
```

## API Endpoints

### Time Tracking

- `POST /api/time/start` - Start a timer
- `POST /api/time/stop` - Stop the current timer
- `GET /api/time/status` - Get current timer status
- `GET /api/time/entries` - Get recent time entries
- `GET /api/time/today` - Get today's summary

### Invoice Generation

- `POST /api/invoice/generate` - Generate an invoice
- `GET /api/invoice/preview` - Preview invoice (TODO)

### Health Check

- `GET /health` - API health status

## Development

### Prerequisites

- Go 1.21+
- Python 3.7+ with the CLI applications installed
- Access to the time tracker and invoice generator CLI tools

### Running the API

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Configure Environment** (Optional):
   ```bash
   # Copy the example environment file
   cp env.example .env
   
   # Edit .env with your specific paths
   nano .env
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

4. **Test the API:**
   ```bash
   curl http://localhost:8080/health
   ```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `PYTHON_EXEC_PATH` | `python3` | Path to Python executable |
| `TIME_TRACKER_PATH` | `../kb-tt-cli` | Path to time tracker CLI |
| `INVOICE_GEN_PATH` | `../kb-invoice-gen-cli` | Path to invoice generator CLI |
| `DATABASE_PATH` | `~/.kb-tt-cli/time_tracker.db` | SQLite database path |

### Example Configuration

```bash
# For conda environment
export PYTHON_EXEC_PATH="/Users/yourusername/anaconda3/envs/your-env/bin/python"

# For custom paths
export TIME_TRACKER_PATH="/path/to/kb-tt-cli"
export INVOICE_GEN_PATH="/path/to/kb-invoice-gen-cli"
export DATABASE_PATH="/path/to/time_tracker.db"

# Start the server
go run main.go
```

## Project Structure

```
kb-freelance-api/
├── main.go                           # Application entry point
├── go.mod                            # Go module definition
├── go.sum                            # Module checksums
├── .gitignore                        # Git ignore rules
├── env.example                       # Environment configuration template
├── run_tests.sh                      # Test runner script
├── run_integration_tests.sh          # Integration test runner
├── README.md                         # This file
├── TESTING.md                        # Testing documentation
├── internal/
│   ├── api/                          # HTTP handlers and server setup
│   │   ├── server.go                 # Gin server configuration
│   │   ├── handlers.go               # API endpoint handlers
│   │   ├── server_test.go            # Server tests
│   │   ├── handlers_test.go          # Handler tests
│   │   └── integration_test.go       # API integration tests
│   ├── config/                       # Configuration management
│   │   ├── config.go                 # Config struct and loading
│   │   └── config_test.go            # Configuration tests
│   └── services/                     # Business logic layer
│       ├── time_tracker.go           # Time tracking service
│       ├── invoice.go                # Invoice generation service
│       ├── time_tracker_test.go      # Time tracker tests
│       ├── invoice_test.go           # Invoice service tests
│       └── integration_test.go       # Service integration tests
└── tests/                            # Test utilities and helpers
```

## Integration with Python CLIs

The API executes the existing Python CLI applications as subprocesses:

- **Time Tracker**: Calls `python3 -m tt.cli` commands
- **Invoice Generator**: Calls `python3 -m src.main` commands

## Testing

The project includes comprehensive testing:

- **Unit Tests**: Mock-based tests for isolated testing
- **Integration Tests**: Real service tests with Python CLI integration
- **Test Coverage**: 61% services, 29% API, 100% config

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests
./run_integration_tests.sh

# Run test suite
./run_tests.sh
```

## Contributing

This is a personal project showcasing Go API development with Python CLI integration. The codebase demonstrates:

- Clean Go project structure with `internal/` packages
- Environment-based configuration
- Comprehensive testing strategy
- Integration with external Python processes
- RESTful API design with Gin framework
