# Architecture Overview

This project follows **Clean Architecture** principles with a clear separation of concerns across different layers.

## Architecture Layers

### 1. Presentation Layer
- **Location**: `internal/handlers/`, `templates/`
- **Purpose**: HTTP request/response handling, HTML rendering
- **Dependencies**: Services layer

### 2. Business Logic Layer  
- **Location**: `internal/services/`
- **Purpose**: Core business logic, use cases, domain rules
- **Dependencies**: Repository layer, Models

### 3. Data Access Layer
- **Location**: `internal/repository/`
- **Purpose**: Data persistence, database operations
- **Dependencies**: Models, Database

### 4. Infrastructure Layer
- **Location**: `internal/database/`, `internal/auth/`, `internal/proxy/`
- **Purpose**: External concerns (database, authentication, proxying)
- **Dependencies**: Configuration

## Directory Structure Explained

### `/cmd/webserver/`
Application entry point. Contains `main.go` which bootstraps the application.

### `/internal/`
Private application code that cannot be imported by other projects.

- **`models/`**: Domain entities and data structures
- **`handlers/`**: HTTP handlers (Controllers in MVC)
- **`services/`**: Business logic and use cases  
- **`repository/`**: Data access layer interfaces and implementations
- **`middleware/`**: HTTP middleware (auth, logging, CORS, etc.)
- **`config/`**: Configuration management
- **`database/`**: Database connection and migration logic
- **`auth/`**: Authentication and authorization utilities
- **`proxy/`**: Reverse proxy functionality for service routing

### `/pkg/`
Public packages that could be imported by other projects.

- **`utils/`**: Common utility functions
- **`logger/`**: Structured logging utilities

### `/templates/`
HTML templates and HTMX components (Views in MVC).

### `/web/static/`
Static assets (CSS, JavaScript, images).

## Data Flow

```
HTTP Request → Middleware → Handler → Service → Repository → Database
                    ↓
HTTP Response ← Template ← Handler ← Service ← Repository ← Database
```

## Dependency Rules

1. **Inward Dependencies Only**: Inner layers should not depend on outer layers
2. **Interface Segregation**: Use interfaces to decouple layers
3. **Dependency Injection**: Inject dependencies rather than creating them
4. **Single Responsibility**: Each package has one clear purpose

## Benefits

- **Testability**: Easy to unit test each layer in isolation
- **Maintainability**: Clear separation makes code easier to understand and modify
- **Flexibility**: Can swap implementations (e.g., SQLite to PostgreSQL) without changing business logic
- **Scalability**: Architecture supports growth and additional features