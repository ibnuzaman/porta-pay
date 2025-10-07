# Booking Service

## Architecture Overview

This service follows Clean Architecture principles with clear separation of concerns:

```
cmd/                    # Application entry point
├── main.go            # Main application

internal/              # Private application code
├── domain/           # Domain layer (business logic)
│   ├── entity/       # Domain entities
│   ├── repository/   # Repository interfaces
│   └── service/      # Service interfaces
├── usecase/          # Application layer (use cases)
├── repository/       # Infrastructure layer (data persistence)
└── delivery/         # Interface layer (HTTP handlers)
    └── http/
        ├── handler/  # HTTP handlers
        ├── middleware/ # HTTP middleware
        └── router/   # Route definitions

migrations/           # Database migrations
api/                 # API documentation
docs/                # Service documentation
```

## Dependencies

- **Router**: Chi v5 for HTTP routing
- **Database**: PostgreSQL with SQLx
- **Logging**: Zerolog
- **Tracing**: OpenTelemetry
- **Config**: Environment variables

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL
- Docker (optional)

### Setup
```bash
# Clone the repository
git clone <repo-url>

# Set environment variables
cp .env.example .env

# Run migrations
make migrate-up

# Start the service
make run
```

### Environment Variables
```
APP_NAME=booking
HTTP_ADDR=:8080
SHUTDOWN_GRACE=10s

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=booking_user
DB_PASSWORD=secret123
DB_NAME=booking_db
DB_SSLMODE=disable

# Observability
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
ENV=dev
```

## Clean Architecture Layers

### 1. Domain Layer (`internal/domain/`)
Contains business logic and rules:
- **Entities**: Core business objects
- **Repository Interfaces**: Data access contracts
- **Service Interfaces**: Business service contracts

### 2. Application Layer (`internal/usecase/`)
Contains application-specific business rules:
- **Use Cases**: Application business logic
- **Dependency injection**: Coordinates between layers

### 3. Infrastructure Layer (`internal/repository/`)
Contains framework and external dependencies:
- **Repository Implementations**: Database access
- **External APIs**: Third-party integrations

### 4. Interface Layer (`internal/delivery/`)
Contains interface adapters:
- **HTTP Handlers**: REST API endpoints
- **Middleware**: HTTP middleware
- **Routers**: Route definitions

## Design Principles

1. **Dependency Inversion**: Dependencies point inward
2. **Single Responsibility**: Each layer has one reason to change
3. **Interface Segregation**: Small, focused interfaces
4. **Separation of Concerns**: Clear boundaries between layers