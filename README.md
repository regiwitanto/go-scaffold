# Go Scaffold Generator

This project is a pure Go backend service built with Echo framework using clean architecture and domain-driven design principles. It aims to replicate the functionality of [Autostrada](https://autostrada.dev/) - a tool that generates REST API scaffolds based on user preferences.

![GitHub language count](https://img.shields.io/github/languages/count/regiwitanto/go-scaffold)
![GitHub top language](https://img.shields.io/github/languages/top/regiwitanto/go-scaffold)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/regiwitanto/go-scaffold)
![GitHub](https://img.shields.io/github/license/regiwitanto/go-scaffold)

## Features

- **REST API Scaffold Generation**: Create complete API application scaffolds with various features and configurations
- **Echo Framework**: High-performance, extensible, minimalist web framework for Go
- **Clean Architecture**: Separation of concerns with layers (domain, application, infrastructure, interfaces)
- **Domain-Driven Design**: Focus on the core domain and domain logic
- **Template Management**: Store and manage templates for different router types
- **Database Support**: Generate scaffolds with different database options (PostgreSQL, MySQL, SQLite)
- **Router Options**: Support for different router implementations (Chi, Echo, Gin, standard library)
- **Configuration Types**: Support for different configuration methods (env vars, flags)
- **Plugin System**: Extensible system to add features like auth, logging, migrations
- **ZIP Generation**: Create downloadable archives of generated code
- **RESTful API**: Complete RESTful API with programmatic access to scaffold generation

## Project Structure

```
go-scaffold/
├── cmd/                      # Main applications of the project
│   └── scaffold/             # The scaffold generator app
│       ├── main.go           # Entry point for the scaffold generator
│       └── main_test.go      # Tests for main.go
├── docs/                     # Documentation
│   └── swagger.go            # OpenAPI/Swagger documentation
├── internal/                 # Private application code
│   ├── application/          # Application layer (use cases)
│   │   └── service/          # Application services
│   ├── domain/               # Domain layer (core business logic)
│   │   ├── model/            # Domain models/entities
│   │   │   ├── scaffold.go   # Scaffold configuration model
│   │   │   └── template.go   # Template model
│   │   ├── repository/       # Repository interfaces
│   │   └── service/          # Domain service interfaces
│   ├── infrastructure/       # Infrastructure layer (implementations)
│   │   └── storage/          # Storage implementations
│   │       ├── scaffold/     # Scaffold storage implementation
│   │       └── template/     # Template storage implementation
│   └── interfaces/           # Interface layer (adapters)
│       ├── api/              # API interfaces
│       │   ├── handler/      # HTTP handlers
│       │   │   ├── generator.go      # Generator handlers
│       │   │   └── template.go       # Template management handlers
│       │   ├── middleware/           # HTTP middleware
│       │   │   ├── logging.go        # Logging middleware
│       │   │   └── ratelimit.go      # Rate limiting middleware
│       │   └── routes/               # Route definitions
│       │       └── routes.go         # API routes
│       └── dto/                      # Data Transfer Objects
│           └── scaffold_request.go   # Scaffold generation request DTO
├── pkg/                              # Public packages
│   ├── logger/                       # Logging utilities
│   │   └── logger.go
│   ├── archiver/                     # ZIP archive creation
│   │   └── zip.go
│   └── errors/                       # Error handling utilities
│       └── errors.go
├── templates/                        # Scaffold templates
│   ├── api/                          # API templates
│   │   ├── chi/                      # Chi router templates
│   │   ├── echo/                     # Echo router templates
│   │   ├── gin/                      # Gin router templates
│   │   └── standard/                 # Standard library templates
│   └── shared/                       # Shared template files
├── test/                             # Test files
│   ├── functional/                   # Functional tests
│   ├── integration/                  # Integration tests
│   └── unit/                         # Unit tests
├── go.mod                            # Go module file
├── go.sum                            # Go dependencies checksums
├── Makefile                          # Build automation
├── .env.example                      # Example environment variables
├── .gitignore                        # Git ignore file
└── README.md                         # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.22+
- Make (optional, for using the Makefile)

### Configuration

Create a `.env` file based on the `.env.example` or set environment variables directly:

```
APP_ENV=development
APP_HTTP_PORT=4000
TEMPLATE_DIR=./templates
```

### Running the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run .

# Or with live reload (requires air: https://github.com/cosmtrek/air)
air
```

The server will start on port 8081 by default. You can change this by setting the PORT environment variable:

```bash
PORT=4000 go run .
```

### Using the Application

The application provides a RESTful API for generating Go API scaffolds. You can interact with it using tools like `curl` or Postman:

```bash
# Generate a scaffold
curl -X POST http://localhost:8081/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "appType": "api",
    "databaseType": "postgresql",
    "routerType": "echo",
    "configType": "env",
    "logFormat": "json",
    "modulePath": "github.com/your-username/your-project",
    "features": ["basic-auth", "sql-migrations"]
  }'

# Download the generated scaffold
curl -o scaffold.zip http://localhost:8081/api/download/GENERATED_ID
```

Available options:

1. API scaffold configuration:
   - Database type: `none`, `postgresql`, `mysql`, `sqlite`
   - Router/framework: `chi`, `echo`, `gin`, `standard`
   - Configuration type: `env`, `flags`
   - Logging format: `json`, `text`

2. Available features:
   - `access-logging`: HTTP request/response logging
   - `admin-makefile`: Administrative Makefile commands
   - `auto-versioning`: Automatic version number generation
   - `basic-auth`: Basic authentication
   - `email-support`: Email sending capabilities
   - `error-notifications`: Error reporting
   - `secure-cookies`: Encrypted cookies
   - `sql-migrations`: Database migrations
   - `user-accounts`: User authentication system
   - `https-support`: HTTPS/TLS support
   - `custom-errors`: Custom error pages
   - And more...

### API Endpoints

- `GET /api/health` - Health check endpoint
- `POST /api/generate` - Generate a scaffold based on provided configuration
- `GET /api/templates` - List available templates
- `GET /api/features` - List available features
- `GET /api/download/:id` - Download a previously generated scaffold
- `GET /api/docs` - API documentation (OpenAPI specification in JSON format)

#### Detailed API Documentation

##### Generate Scaffold

```
POST /api/generate
```

Request body:

```json
{
  "appType": "api",
  "databaseType": "postgresql",
  "routerType": "echo",
  "configType": "env",
  "logFormat": "json",
  "modulePath": "github.com/username/project",
  "features": [
    "basic-auth",
    "sql-migrations",
    "https-support"
  ]
}
```

Response:

```json
{
  "id": "scaffold-123456",
  "options": {
    "appType": "api",
    "databaseType": "postgresql",
    "routerType": "echo",
    "configType": "env",
    "logFormat": "json",
    "modulePath": "github.com/username/project",
    "features": ["basic-auth", "sql-migrations", "https-support"]
  },
  "createdAt": "2025-07-09T12:34:56Z",
  "filePath": "/tmp/go-scaffold/scaffold-123456.zip",
  "size": 123456
}
```

##### List Templates

```
GET /api/templates
```

Response:

```json
[
  {
    "id": "api-echo",
    "name": "API with Echo",
    "description": "RESTful API using Echo framework",
    "type": "api",
    "path": "/templates/api/echo"
  },
  {
    "id": "api-chi",
    "name": "API with Chi",
    "description": "RESTful API using Chi router",
    "type": "api",
    "path": "/templates/api/chi"
  }
]
```

##### List Features

```
GET /api/features
```

Response:

```json
[
  {
    "id": "basic-auth",
    "name": "Basic Authentication",
    "description": "Adds HTTP Basic Authentication to the API",
    "isPremium": false
  },
  {
    "id": "sql-migrations",
    "name": "SQL Migrations",
    "description": "Adds database migration support",
    "isPremium": false
  }
]
```

##### Download Scaffold

```
GET /api/download/:id
```

Returns a ZIP file containing the generated scaffold.

### Client Libraries

We provide official client libraries to interact with the Go Scaffold Generator API:

- [go-scaffold-client-go](https://github.com/regiwitanto/go-scaffold-client-go) - Go client
- [go-scaffold-client-js](https://github.com/regiwitanto/go-scaffold-client-js) - JavaScript/TypeScript client

These libraries make it easy to integrate Go Scaffold Generator with your development workflow or tools.

### Docker Deployment

You can also run this service using Docker:

```bash
# Build the Docker image
docker build -t go-scaffold .

# Run the container
docker run -p 8081:8081 -e PORT=8081 go-scaffold
```

For production deployment, you might want to use Docker Compose:

```yaml
# docker-compose.yml
version: '3'
services:
  go-scaffold:
    build: .
    ports:
      - "8081:8081"
    environment:
      - PORT=8081
      - APP_ENV=production
      - TEMPLATE_DIR=/app/templates
    volumes:
      - ./templates:/app/templates
    restart: always
```

Run with `docker-compose up -d`.

## Development

### Core Components

1. **Template Engine**: Responsible for parsing and processing templates
2. **Scaffold Generator**: Creates REST API scaffolds based on user preferences
3. **API**: RESTful API for programmatic scaffold generation
4. **ZIP Creator**: Creates downloadable archives of generated code

### Architecture Overview

This project follows Clean Architecture with Domain-Driven Design principles:

1. **Domain Layer**: 
   - Core entities and business rules
   - Independent of external frameworks and tools

2. **Application Layer**:
   - Use cases that orchestrate the flow of data
   - Scaffold generation logic

3. **Infrastructure Layer**:
   - Template storage (filesystem)
   - ZIP file generation

4. **Interface Layer**:
   - RESTful API
   - DTOs for data transformation

### Extension Points

The system is designed to be extensible through:

1. **Template System**: Add new templates for different router types and frameworks
2. **Feature Plugins**: Add new features that can be included in the generated code
3. **Output Formats**: Support different output formats beyond ZIP files

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Autostrada](https://autostrada.dev/) for inspiration
- [Echo Framework](https://echo.labstack.com/) for the web framework
- [text/template](https://golang.org/pkg/text/template/) for template processing
