# Web Server Architecture

This document explains the clean architecture implementation of the Homer Dashboard Server's web interface.

## Architecture Overview

The web server follows **Clean Architecture** principles with clear separation of concerns:

```
cmd/webserver/main.go          # Application entry point
├── internal/server/           # HTTP server setup and configuration
├── internal/handlers/         # HTTP request handlers (Controllers)
├── internal/services/         # Business logic layer
├── internal/models/           # Data structures and entities
└── templates/                 # HTML templates (Views)
```

## Layer Responsibilities

### 1. **Entry Point** (`cmd/webserver/main.go`)
- Application bootstrap
- Server initialization
- Minimal, focused on startup

### 2. **Server Layer** (`internal/server/`)
- HTTP server configuration
- Route setup and middleware
- Template loading and management
- Dependency injection

### 3. **Handlers Layer** (`internal/handlers/`)
- HTTP request/response handling
- Input validation and sanitization
- Template rendering
- Error handling

### 4. **Services Layer** (`internal/services/`)
- Business logic implementation
- File system operations
- Data processing and transformation
- Security validations

### 5. **Models Layer** (`internal/models/`)
- Data structures and entities
- Template data models
- Domain objects

### 6. **Templates** (`templates/`)
- HTML presentation layer
- Reusable layout components
- HTMX integration

## File Structure

### Templates
- **`layout.html`**: Base layout with CSS and common structure
- **`index.html`**: Main page content
- **`browse.html`**: Directory browsing interface
- **`file.html`**: File content display

### Go Packages
- **`models/file.go`**: File system data structures
- **`services/file_service.go`**: File operations business logic
- **`handlers/file_handler.go`**: HTTP request handlers
- **`server/server.go`**: Server setup and configuration

## Benefits of This Architecture

### **Separation of Concerns**
- Each layer has a single responsibility
- Easy to test individual components
- Clear boundaries between layers

### **Maintainability**
- HTML is separate from Go code
- Business logic is isolated in services
- Easy to modify templates without touching Go code

### **Testability**
- Each layer can be unit tested independently
- Mock interfaces for testing
- Clear dependencies

### **Scalability**
- Easy to add new features
- Can swap implementations (e.g., different file systems)
- Supports middleware and additional handlers

### **Security**
- Input validation in handlers
- Security logic centralized in services
- Template auto-escaping

## Request Flow

```
HTTP Request → Server → Handler → Service → Model
                ↓
HTTP Response ← Template ← Handler ← Service ← Model
```

1. **Request arrives** at the server
2. **Router** directs to appropriate handler
3. **Handler** validates input and calls service
4. **Service** performs business logic
5. **Model** structures the data
6. **Template** renders the response
7. **Handler** sends response to client

## Template System

Uses Go's built-in `html/template` package with:
- **Template inheritance** via `layout.html`
- **Automatic HTML escaping** for security
- **Reusable components** with `{{define}}` blocks
- **Data binding** with struct fields

## HTMX Integration

- **Dynamic content loading** without page refreshes
- **Progressive enhancement** - works without JavaScript
- **Server-side rendering** with client-side interactivity
- **Minimal JavaScript footprint**

## Security Features

- **Path traversal protection** in file service
- **Input sanitization** in handlers
- **HTML escaping** in templates
- **Content-Type headers** for proper encoding
- **Error handling** without information leakage

This architecture provides a solid foundation for building more complex features while maintaining clean, testable, and maintainable code.