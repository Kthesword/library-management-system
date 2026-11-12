# Console-Based Library Management System (Go)

A simple, modular console-based library management system written in Go to demonstrate structs, interfaces, methods, slices, maps, and Docker containerization.

## Project Structure

```
library_management/
├── main.go                     # Entry point of the application
├── controllers/
│   ├── library_controller.go   # Handles console input and delegates to services
│   └── library_controller_test.go
├── models/
│   ├── book.go                 # Defines Book struct and statuses
│   └── member.go               # Defines Member struct
├── services/
│   ├── library_service.go      # LibraryManager interface and Library implementation
│   └── library_service_test.go
├── docs/
│   └── documentation.md        # Comprehensive system documentation
├── Dockerfile                  # Docker container build definition
├── docker-compose.yml          # Docker Compose configuration
└── go.mod                      # Module definition
```

## Quick Start

### 1. Run with Go Locally

```bash
cd library_management
go run main.go
```

### 2. Run with Docker

Build the Docker image:
```bash
docker build -t library-management .
```

Run the container interactively:
```bash
docker run -it --rm library-management
```

Or using Docker Compose:
```bash
docker compose run --rm library-app
```

### 3. Run Unit Tests

```bash
cd library_management
go test -v ./...
```

For full details, see [docs/documentation.md](file:///c:/Users/jafer/Desktop/Go-Library-management-system/library_management/docs/documentation.md).

