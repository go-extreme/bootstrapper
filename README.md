# Bootstrapper

A lightweight Go package for managing application initialization order with dependency injection support.

## Features

- Register components with custom initialization order
- Thread-safe shared state management
- Automatic lifecycle management (Booting/Booted phases)
- Predefined order constants for common components

## Usage

### Basic Registration

```go
package main

type Database struct{}

func (d *Database) Booting() {
    // Initialize database connection
}

func (d *Database) Booted() {
    // Run migrations or post-init tasks
}

func main() {
    db := &Database{}
    bootstrapper.Register(db, bootstrapper.OrderDB)
    
    bootstrapper.StartBooting()
    bootstrapper.StartBooted()
}
```

### Custom Order

```go
const CustomOrder = 50

service := &MyService{}
bootstrapper.Register(service, CustomOrder)
```

### Shared State

```go
// Store shared data
bootstrapper.Share("config", myConfig)

// Retrieve shared data
config := bootstrapper.Get("config")
```

## Predefined Orders

- `OrderConfig` (1) - Configuration
- `OrderDB` (2) - Database
- `OrderCache` (3) - Cache
- `OrderQueue` (4) - Queue
- `WorkerPool` (5) - Worker Pool
- `OrderRepository` (6) - Repository
- `OrderService` (7) - Service
- `OrderController` (8) - Controller

## Testing

```bash
go test -v
```

## License

MIT
