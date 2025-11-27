# Testing Guide

## Unit Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific service
make test-service SERVICE=master-data

# Run tests with race detector
make test-race
```

### Test Structure

Tests are located alongside the code they test:

```
internal/
  master-data/
    service/
      service.go
      service_test.go
  leave/
    service/
      service.go
      service_test.go
```

### Writing Tests

1. **Use table-driven tests** for multiple test cases:

```go
func TestCreateStudyProgram(t *testing.T) {
    tests := []struct {
        name    string
        req     CreateStudyProgramRequest
        wantErr bool
    }{
        {
            name: "success",
            req: CreateStudyProgramRequest{
                Code: "TI",
                Name: "Teknik Informatika",
            },
            wantErr: false,
        },
        {
            name: "duplicate code",
            req: CreateStudyProgramRequest{
                Code: "TI",
                Name: "Duplicate",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

2. **Use mocks** for dependencies:

```go
type MockRepository struct {
    // Mock implementation
}
```

3. **Test edge cases**:
   - Invalid input
   - Not found errors
   - Duplicate entries
   - Boundary conditions

### Coverage Goals

- Minimum coverage: **70%**
- Target coverage: **80%+**
- Critical paths: **90%+**

### Running Coverage Report

```bash
# Generate coverage report
make test-coverage

# View HTML report
open coverage.html
```

## Integration Testing

### Setup Test Database

```bash
# Start test database
docker-compose -f deployments/docker-compose/docker-compose.test.yml up -d

# Run migrations
make migrate-up
```

### API Testing

Use Postman collection:

```bash
# Import collection
postman/UNSRI_Backend_API.postman_collection.json
```

## CI/CD Testing

Tests run automatically on:

- Push to `main` or `develop` branches
- Pull requests
- Before deployment

See `.github/workflows/ci.yml` for details.
