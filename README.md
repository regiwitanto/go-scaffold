# Go Scaffold Generator

A pure Go backend service that generates REST API scaffolds based on user preferences, built with Echo framework.

## Features

- Generate complete API scaffolds with Echo, Chi, Gin, or standard library
- Database support (PostgreSQL, MySQL, SQLite)
- Clean Architecture & Domain-Driven Design
- Plugin system with features like auth, logging, migrations
- Configuration options (env vars, flags)
- RESTful API for programmatic access

## Quick Start

### Prerequisites
- Go 1.22+
- Make (optional)

### Install & Run

```bash
# Clone and build
git clone https://github.com/regiwitanto/go-scaffold.git
cd go-scaffold
go mod tidy
make build  # or go build -o build/go-scaffold ./cmd/scaffold

# Run the server
go run ./cmd/scaffold  # default port 8081
```

### Configuration
Create a `.env` file or set environment variables:
```
APP_ENV=development
PORT=8081
TEMPLATE_DIR=./templates
```

## API Usage

```bash
# Generate a scaffold
curl -X POST http://localhost:8081/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "appType": "api",
    "databaseType": "postgresql",
    "routerType": "echo",
    "configType": "env",
    "modulePath": "github.com/username/project",
    "features": ["basic-auth", "sql-migrations"]
  }'

# Download the scaffold
curl -o project.zip http://localhost:8081/api/download/SCAFFOLD_ID
```

### API Endpoints

- `GET /api/health` - Health check
- `POST /api/generate` - Generate scaffold
- `GET /api/templates` - List templates
- `GET /api/features` - List features
- `GET /api/download/:id` - Download scaffold
- `GET /api/docs` - API documentation

### Docker

```bash
# Run with Docker
docker build -t go-scaffold .
docker run -p 8081:8081 -e PORT=8081 go-scaffold

# Or with Docker Compose
docker-compose up -d
```

## Development

### Testing

```bash
# Run tests
make test            # all tests
make test-unit       # unit tests only
make test-cover      # with coverage
```

### Architecture

The project follows Clean Architecture with:
- **Domain Layer**: Core business logic
- **Application Layer**: Use cases
- **Infrastructure Layer**: Storage implementations
- **Interface Layer**: API and controllers

## License & Contributing

- Licensed under the MIT License
- Contributions welcome through issues and pull requests

## Acknowledgments

- [Autostrada](https://autostrada.dev/) for inspiration
- [Echo Framework](https://echo.labstack.com/) for the web framework
