# Test Operations

This directory contains test configurations and scripts for the go-metadata-removal-toolkit project.

## Test Types

### Unit Tests
```bash
# Run all unit tests
go test ./src/...

# Run with coverage
go test -cover ./src/...

# Generate coverage report
go test -coverprofile=coverage.out ./src/...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./src/...
```

### Benchmark Tests
```bash
# Run benchmarks
go test -bench=. -benchmem ./src/...

# Run specific benchmark
go test -bench=BenchmarkProcessImage ./src/processor
```

### End-to-End Tests
```bash
# Build and test the binary
go build -o metadata-remover ./src
./metadata-remover -path test/fixtures -preview
```

## Test Coverage

We maintain a 1:1 source-to-test file ratio:
- Every `.go` file has a corresponding `_test.go` file
- Minimum coverage target: 80%
- Critical paths require 100% coverage

## Test Data

Test fixtures are located in:
- `test/fixtures/` - Sample files for testing
- `test/golden/` - Expected output for comparison

## Continuous Integration

Tests are automatically run on:
- Every pull request
- Every commit to main branch
- Nightly scheduled runs

## Performance Testing

```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=. ./src/...
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=. ./src/...
go tool pprof mem.prof
```

## Test Guidelines

1. **Naming**: Use descriptive test names (TestFunctionName_Scenario_ExpectedResult)
2. **Table Tests**: Prefer table-driven tests for multiple scenarios
3. **Isolation**: Tests should not depend on external resources
4. **Cleanup**: Always clean up test artifacts
5. **Assertions**: Use clear assertion messages

## Running Tests in Container

```bash
# Run tests in development container
podman run --rm -v $(pwd):/workspace metadata-remover:dev \
    go test -v ./src/...
```

## Mocking

For mocking external dependencies:
```go
//go:generate mockgen -source=interface.go -destination=mock_interface.go
```

## Test Reports

Test results are generated in:
- `coverage.html` - HTML coverage report
- `test-results.xml` - JUnit format for CI
- `benchmark-results.txt` - Benchmark comparisons
