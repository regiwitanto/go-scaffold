# Contributing to Go Scaffold Generator

Thank you for your interest in contributing to Go Scaffold Generator! This document provides guidelines and instructions to help you get started.

## Development Setup

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/go-scaffold.git
cd go-scaffold/echo-scaffold
```

2. **Install dependencies**

```bash
go mod tidy
```

3. **Run the application**

```bash
make run
# or
go run .
```

4. **Run tests**

```bash
make test
# or
go test ./...
```

## Project Structure

The project follows clean architecture and domain-driven design principles:

- `cmd/`: Application entry points
- `internal/`: Internal packages
  - `domain/`: Business domain layer (models, interfaces)
  - `application/`: Application services
  - `infrastructure/`: External implementations (storage, template engines)
  - `interfaces/`: User interfaces (API, web)
- `templates/`: Scaffold templates
- `assets/`: Static assets

## Adding New Templates

To add a new template type:

1. Create a new directory in the `templates/` directory with the appropriate structure
2. Register the template in the template repository implementation
3. Add tests for your new template

Example for adding a new API template with Fiber framework:

```bash
mkdir -p templates/api/fiber
# Add template files
```

## Adding New Features

To add a new feature:

1. Implement the feature in the appropriate package
2. Register the feature in the generator service
3. Add tests for your new feature

## Code Guidelines

- Follow standard Go coding conventions
- Add tests for new code
- Document public functions and packages
- Use meaningful variable and function names
- Keep functions short and focused

## Submitting Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Add tests for your changes
5. Run tests (`go test ./...`)
6. Commit your changes (`git commit -am 'Add new feature'`)
7. Push to the branch (`git push origin feature/my-feature`)
8. Create a Pull Request

## License

By contributing, you agree that your contributions will be licensed under the project's MIT License.
