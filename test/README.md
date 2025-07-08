# Go-Scaffold Test Suite

This directory contains the test suite for the Go-Scaffold project. The tests are organized into three main categories:

## Test Structure

### 1. Unit Tests (`/test/unit`)
Unit tests focus on testing individual components in isolation. These tests are fast, don't require external dependencies, and are used to validate the correctness of specific functions or methods.

Example unit tests include:
- Testing domain model validation
- Testing template parsing logic
- Testing API handlers with mock services

### 2. Integration Tests (`/test/integration`)
Integration tests focus on testing the interaction between components. These tests validate that different parts of the system work together correctly.

Example integration tests include:
- Testing template rendering with actual template files
- Testing database interactions with a test database
- Testing API endpoints with actual services

### 3. Functional Tests (`/test/functional`)
Functional tests focus on testing the entire system from an external perspective. These tests validate that the system meets its requirements and works correctly from a user's perspective.

Example functional tests include:
- Testing the CLI tool to generate scaffolds
- Testing the API server to generate scaffolds
- Testing the generated scaffolds to ensure they compile and run correctly

## Running Tests

You can use the following Make targets to run the tests:

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run functional tests only
make test-functional

# Run tests with coverage
make test-cover
```

## Writing New Tests

When writing new tests:

1. Identify the category your test belongs to (unit, integration, functional)
2. Place the test in the appropriate directory
3. Name the test file with the `_test.go` suffix
4. Use the appropriate mocks from the `/test/mocks` directory
5. Use the test utilities from the `/test/testutil` package

## Mock Objects

Mock objects for testing are located in the `/test/mocks` directory. These mocks implement the interfaces defined in the domain layer, allowing for isolated testing of components.

## Test Utilities

Common test utilities are located in the `/test/testutil` package. These utilities provide common functionality for testing, such as finding the project root, asserting file existence, etc.
