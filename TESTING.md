# Testing Documentation

This document describes the testing strategy and coverage for the kb-freelance-api Go application.

## Test Structure

The tests are organized by package and functionality:

```
kb-freelance-api/
├── main_test.go                    # Main package tests
├── internal/
│   ├── config/
│   │   └── config_test.go         # Configuration tests
│   ├── services/
│   │   ├── time_tracker_test.go   # Time tracker service tests
│   │   └── invoice_test.go        # Invoice service tests
│   └── api/
│       ├── handlers_test.go       # API handler tests
│       └── server_test.go         # Server setup tests
└── run_tests.sh                   # Test runner script
```

## Test Coverage

### Configuration Tests (`internal/config`)
- **Coverage**: 100%
- **Tests**:
  - `TestLoad()` - Tests configuration loading and path resolution
  - `TestGetEnv()` - Tests environment variable handling
  - `TestConfigPaths()` - Tests path validation and absolute path generation

### Services Tests (`internal/services`)
- **Coverage**: 3.9%
- **Time Tracker Service Tests**:
  - `TestNewTimeTrackerService()` - Service initialization
  - `TestTimerStatus()` - Timer status struct validation
  - `TestTimeEntry()` - Time entry struct validation
  - `TestTodaySummary()` - Today summary struct validation
  - `TestBreakdown()` - Breakdown struct validation
  - `TestContains()` - Helper function testing
  - `TestParseTodaySummary()` - Regex parsing logic testing

- **Invoice Service Tests**:
  - `TestNewInvoiceService()` - Service initialization
  - `TestInvoiceLineItem()` - Line item struct validation
  - `TestInvoiceLineItemValidation()` - Input validation testing
  - `TestInvoiceLineItemCalculations()` - Calculation logic testing

### API Tests (`internal/api`)
- **Coverage**: 1.4%
- **Handler Tests**:
  - `TestHealthEndpoint()` - Health check endpoint
  - `TestStartTimerEndpoint()` - Start timer with mock service
  - `TestStopTimerEndpoint()` - Stop timer with mock service
  - `TestGetStatusEndpoint()` - Get timer status with mock service
  - `TestGetTodaySummaryEndpoint()` - Get today's summary with mock service
  - `TestGenerateInvoiceEndpoint()` - Generate invoice with mock service
  - `TestInvalidJSONRequest()` - Error handling for invalid JSON

- **Server Tests**:
  - `TestNewServer()` - Server initialization
  - `TestServerRoutes()` - Route registration validation
  - `TestCORSConfiguration()` - CORS middleware testing
  - `TestServerMiddleware()` - Middleware functionality testing

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Tests with Verbose Output
```bash
go test -v ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Tests with Detailed Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

### Use the Test Runner Script
```bash
./run_tests.sh
```

## Test Strategy

### Unit Tests
- **Mocking**: Uses `github.com/stretchr/testify/mock` for service mocking
- **Assertions**: Uses `github.com/stretchr/testify/assert` for test assertions
- **HTTP Testing**: Uses `net/http/httptest` for HTTP endpoint testing

### Test Categories

1. **Configuration Tests**: Test environment setup and path resolution
2. **Service Tests**: Test business logic and data structures
3. **Handler Tests**: Test HTTP request/response handling with mocked services
4. **Integration Tests**: Test server setup and middleware configuration

### Mock Strategy

The tests use mock services to isolate the API layer from the service layer:

```go
type MockTimeTrackerService struct {
    mock.Mock
}

func (m *MockTimeTrackerService) StartTimer(client, project, description string) (map[string]interface{}, error) {
    args := m.Called(client, project, description)
    return args.Get(0).(map[string]interface{}), args.Error(1)
}
```

## Test Data

### Sample Output Testing
The tests include sample output from the Python CLI to test parsing logic:

```go
sampleOutput := `╭───────────────────────────── 📊 Today's Summary ─────────────────────────────╮
│ Total Time: 2.5 hours                                                        │
│                                                                              │
│ Breakdown:                                                                   │
│   • Client A/Project A: 1.5h                                                 │
│   • Client B/Project B: 1.0h                                                 │
│                                                                              │
╰──────────────────────────────────────────────────────────────────────────────╯`
```

## Coverage Goals

- **Current Coverage**: 
  - Config: 100%
  - Services: 3.9%
  - API: 1.4%
  - Overall: Low due to integration with external Python processes

- **Target Coverage**: 
  - Unit testable code: 80%+
  - Integration points: Tested with mocks

## Future Improvements

1. **Integration Tests**: Add tests that actually call the Python CLI processes
2. **Performance Tests**: Add benchmarks for critical paths
3. **Error Scenarios**: Add more error handling tests
4. **Edge Cases**: Test boundary conditions and edge cases
5. **Concurrency Tests**: Test concurrent API calls

## Test Dependencies

- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/mock` - Mocking
- `github.com/gin-gonic/gin` - HTTP framework testing
- `net/http/httptest` - HTTP testing utilities

## Continuous Integration

The tests are designed to run in CI/CD pipelines and should:
- Pass without external dependencies (Python CLI processes)
- Complete in under 30 seconds
- Provide clear failure messages
- Generate coverage reports
