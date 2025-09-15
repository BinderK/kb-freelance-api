# KB Freelance API

A Go-based REST API that integrates with the existing Python CLI applications for time tracking and invoice generation.

## Features

- **Time Tracking**: Start/stop timers, get status, view entries
- **Invoice Generation**: Generate invoices from time entries
- **REST API**: Clean RESTful interface for frontend integration
- **CORS Support**: Ready for React frontend integration

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
├── main.go                    # Application entry point
├── internal/
│   ├── api/                   # HTTP handlers and server setup
│   │   ├── server.go         # Gin server configuration
│   │   └── handlers.go       # API endpoint handlers
│   ├── config/               # Configuration management
│   │   └── config.go         # Config struct and loading
│   └── services/             # Business logic layer
│       ├── time_tracker.go   # Time tracking service
│       └── invoice.go        # Invoice generation service
└── README.md                 # This file
```

## Integration with Python CLIs

The API executes the existing Python CLI applications as subprocesses:

- **Time Tracker**: Calls `python3 -m tt.cli` commands
- **Invoice Generator**: Calls `python3 -m src.main` commands

## Future Enhancements

- [ ] Add JSON output to Python CLIs for better parsing
- [ ] Implement database direct access for better performance
- [ ] Add authentication and authorization
- [ ] Add comprehensive logging and monitoring
- [ ] Add API documentation with Swagger
- [ ] Add unit and integration tests
