# Go-Scaffold Template Guide

This document explains the templates and features available in Go-Scaffold, and how to maintain and extend them.

## Template Organization

Go-Scaffold templates are organized by application type and router framework:

```
templates/
├── api/
│   ├── chi/
│   ├── echo/
│   ├── gin/
│   └── standard/
├── shared/
│   ├── auth/
│   ├── db/
│   └── email/
└── webapp/
    ├── chi/
    ├── echo/
    ├── gin/
    └── standard/
```

## Feature Implementation

Each feature is implemented as conditional blocks in templates using Go's template syntax. Features can be checked with the `HasFeature` function:

```go
{{if (call .HasFeature "feature-name")}}
// Feature-specific code here
{{end}}
```

### Core Features

| Feature ID | Description | File Locations |
|------------|-------------|----------------|
| access-logging | HTTP request logging | middleware/, main.go |
| admin-makefile | Development tasks | Makefile |
| automatic-versioning | VCS-based version | version/ |
| basic-auth | HTTP authentication | middleware/auth.go |
| email | SMTP email support | email/mailer.go |
| error-notifications | Admin alerts | handlers/, middleware/ |
| gitignore | Version control ignore | .gitignore |
| live-reload | Auto rebuild and restart | Makefile, scripts/ |
| secure-cookies | Cookie encryption | cookies/ |
| sql-migrations | Database migrations | database/db.go, migrations/ |

### Premium Features

| Feature ID | Description | File Locations |
|------------|-------------|----------------|
| automatic-https | Let's Encrypt TLS | server/ |
| custom-error-pages | Custom HTML errors | templates/, handlers/ |
| user-accounts | User management | handlers/, models/ |

## Maintaining Templates

When modifying templates, ensure:

1. **Feature Consistency**: All templates should handle features consistently
2. **Database Support**: All templates must support all database types (PostgreSQL, MySQL, SQLite)
3. **Config Options**: Both environment variables and flags should be supported
4. **Error Handling**: Proper error handling and logging

## Adding New Features

To add a new feature:

1. Add it to the feature list in `internal/application/service/generator_service.go`
2. Implement the feature in shared templates if possible
3. Add conditional blocks to all relevant templates
4. Update this documentation
5. Add tests to verify the feature

## Testing Templates

Use the provided test scripts to validate templates:

1. Unit tests: `go test ./internal/infrastructure/storage/template/`
2. Integration tests: `./test/test-combinations.sh`

## Common Troubleshooting

- **Syntax Errors**: Check for mismatched template tags (`{{` and `}}`)
- **Missing Imports**: Ensure conditional imports are correctly implemented
- **Database Errors**: Verify DSN construction for all database types
- **Configuration Errors**: Check that all config options are properly parsed
