# Go Scaffold Generator

This project is a Go backend service built with Echo framework using clean architecture and domain-driven design principles. It aims to replicate the functionality of [Autostrada](https://autostrada.dev/) - a tool that generates application scaffolds based on user preferences.

## Features

- **Scaffold Generation**: Create complete application scaffolds with various features and configurations
- **Echo Framework**: High-performance, extensible, minimalist web framework for Go
- **Clean Architecture**: Separation of concerns with layers (domain, application, infrastructure, interfaces)
- **Domain-Driven Design**: Focus on the core domain and domain logic
- **Template Management**: Store and manage templates for different application types
- **Database Support**: Generate scaffolds with different database options (PostgreSQL, MySQL, SQLite)
- **Router Options**: Support for different router implementations (Chi, Echo, etc.)
- **Configuration Types**: Support for different configuration methods (env vars, flags)
- **Plugin System**: Extensible system to add features like auth, logging, migrations
- **ZIP Generation**: Create downloadable archives of generated code
- **Web Interface**: User-friendly web interface to configure and generate scaffolds

## Project Structure

```
echo-scaffold/
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
│       └── api/              # API interfaces
│           ├── handler/      # HTTP handlers
│           └── routes/       # Route definitions
├── templates/                # Scaffold templates
│   ├── api/                  # API templates
│   │   ├── chi/              # Chi router templates
│   │   ├── echo/             # Echo router templates
│   │   ├── gin/              # Gin router templates
│   │   └── standard/         # Standard library templates
│   └── webapp/               # Web application templates
├── assets/                   # Static assets
├── Makefile                  # Build automation
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── .env                      # Environment variables
└── .env.example              # Example environment configuration
│       ├── api/                      # API handlers
│       │   ├── handler/              # HTTP handlers
│       │   │   ├── generator.go      # Generator handlers
│       │   │   └── template.go       # Template management handlers
│       │   ├── middleware/           # HTTP middleware
│       │   │   ├── logging.go        # Logging middleware
│       │   │   └── ratelimit.go      # Rate limiting middleware
│       │   └── routes/               # Route definitions
│       │       └── routes.go         # API routes
│       ├── dto/                      # Data Transfer Objects
│       │   └── scaffold_request.go   # Scaffold generation request DTO
│       └── web/                      # Web interface
│           └── templates/            # Web UI templates
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
│   │   └── standard/                 # Standard library templates
│   └── webapp/                       # Web application templates
│       ├── standard/                 # Standard library templates
│       └── templates/                # UI templates
├── assets/                           # Static assets for web UI
│   ├── css/                          # CSS styles
│   ├── js/                           # JavaScript files
│   └── images/                       # Images
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

Once running, you can access the web interface at http://localhost:4000 which allows you to:

1. Choose an application type (API or Web Application)
2. Customize the application settings:
   - Database type (None, PostgreSQL, MySQL, SQLite)
   - Router/framework (Chi, Echo, standard library)
   - Configuration type (environment variables, command-line flags)
   - Logging format (JSON, plain text)
3. Select additional features:
   - Access logging
   - Admin Makefile
   - Automatic versioning
   - Basic auth
   - Email support
   - Error notifications
   - Secure cookies
   - SQL migrations
   - User accounts and authentication
   - HTTPS support
   - Custom error pages
   - And more...

After configuring, the application will generate a ZIP file containing the scaffold code that can be downloaded.

### API Endpoints

- `GET /api/health` - Health check endpoint
- `POST /api/generate` - Generate a scaffold based on provided configuration
- `GET /api/templates` - List available templates
- `GET /api/features` - List available features
- `GET /api/download/:id` - Download a previously generated scaffold
- `GET /api/docs` - API documentation (OpenAPI specification in JSON format)
- `GET /api-docs` - Interactive API documentation using Swagger UI

## Development

### Core Components

1. **Template Engine**: Responsible for parsing and processing templates
2. **Scaffold Generator**: Creates application scaffolds based on user preferences
3. **Web Interface**: User-friendly interface for configuring and generating scaffolds
4. **API**: RESTful API for programmatic scaffold generation
5. **ZIP Creator**: Creates downloadable archives of generated code

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
   - Web UI
   - RESTful API
   - DTOs for data transformation

### Extension Points

The system is designed to be extensible through:

1. **Template System**: Add new templates for different frameworks or application types
2. **Feature Plugins**: Add new features that can be included in the generated code
3. **Output Formats**: Support different output formats beyond ZIP files

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Autostrada](https://autostrada.dev/) for inspiration
- [Echo Framework](https://echo.labstack.com/) for the web framework
- [text/template](https://golang.org/pkg/text/template/) for template processing
