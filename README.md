# Zero - Minimal HTTP Server Library

A simple, lean, and blazingly fast HTTP server library built with pure Go using only native libraries. Designed for maximum performance, security, and minimal resource footprint. Zero provides a lightweight foundation for building high-performance HTTP servers with comprehensive middleware support and graceful shutdown capabilities.

## 🚀 Features

- **Zero Dependencies**: Uses only Go's standard library - no external packages
- **Lightning Fast**: Optimized for high-performance request handling
- **Memory Efficient**: Minimal memory footprint and garbage collection pressure
- **Security First**: Built with security best practices and input validation
- **Simple & Clean**: Easy to understand, maintain, and extend
- **Production Ready**: Robust error handling, logging, and graceful shutdown
- **Middleware Support**: Built-in CORS, logging, and recovery middleware
- **Comprehensive Documentation**: Fully documented with Go standard comments
- **Type-Safe**: Strong typing for request parameters and responses
- **Context Integration**: Full support for request context and tracing

## 🎯 Why Zero?

- **Simplicity**: No complex frameworks or abstractions
- **Performance**: Direct use of Go's `net/http` for maximum speed
- **Security**: Reduced attack surface with minimal dependencies
- **Reliability**: Fewer moving parts means fewer failure points
- **Maintainability**: Clean, readable code that's easy to debug
- **Flexibility**: Easy to customize and extend for your needs

## 📋 Requirements

- Go 1.24.7 or higher
- No external dependencies required

## 🛠️ Installation

```bash
# Add to your project
go get github.com/aleal/zero

# Or clone the repository
git clone https://github.com/aleal/zero.git
cd zero

# Run the example server
go run examples/server/cmd/main.go
```

## 🚀 Quick Start

```bash
# Run the example server
go run examples/server/cmd/main.go

# Or build and run the example
go build -o server examples/server/cmd/main.go && ./server
```

The example server will start on `http://localhost:8000` by default.

### Using the Library

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/aleal/zero/pkg/server"
    "github.com/aleal/zero/pkg/request"
    "github.com/aleal/zero/pkg/response"
)

func main() {
    ctx := context.Background()
    
    // Create a new server instance
    server := server.NewServer(ctx)
    
    // Add your routes
    server.Get("/hello", func(rctx context.Context, w http.ResponseWriter, r *http.Request) {
        response.WriteJSON(w, http.StatusOK, map[string]string{
            "message": "Hello from Zero!",
        })
    })
    
    // Start the server
    server.Start()
}
```

## 📖 Usage

### Basic Usage

```bash
# Run the example server with default settings
go run examples/server/cmd/main.go

# Run with custom port
go run examples/server/cmd/main.go -port 3000

# Run with custom host
go run examples/server/cmd/main.go -host 0.0.0.0 -port 8000
```

### Environment Variables

The server supports configuration through environment variables:

- `ZERO_HOST`: Server host (default: localhost)
- `ZERO_PORT`: Server port (default: 8000)
- `ZERO_READ_TIMEOUT`: Read timeout duration (default: 5s)
- `ZERO_WRITE_TIMEOUT`: Write timeout duration (default: 15s)
- `ZERO_IDLE_TIMEOUT`: Idle timeout duration (default: 60s)
- `ZERO_RATE_LIMIT`: Rate limit per second (default: 100)
- `ZERO_ENABLE_LOGGING`: Enable request logging (default: true)
- `ZERO_ENABLE_CORS`: Enable CORS middleware (default: true)
- `ZERO_ENABLE_RECOVERY`: Enable panic recovery (default: true)
- `ZERO_LOG_LEVEL`: Log level (DEBUG, INFO, WARNING, ERROR, FATAL, PANIC)

## 🏗️ Architecture

```
zero/
├── pkg/                   # Public packages
│   ├── server/            # Core server implementation
│   ├── config/            # Configuration management
│   ├── log/               # Structured logging
│   ├── middlewares/       # HTTP middleware
│   ├── request/           # Request utilities
│   ├── response/          # Response utilities
│   └── uuid/              # UUID generation
├── internal/              # Internal packages
│   ├── context/           # Context key definitions
│   └── handlers/          # Built-in HTTP handlers
├── examples/              # Runnable examples
│   └── server/            # Example server implementation
│       ├── cmd/           # Example server entry point
│       └── static/        # Static files
└── README.md              # This file
```

## 🔧 Configuration

The library provides flexible configuration through:

- **Programmatic Configuration**: Direct configuration via `config.Config`
- **Environment Variables**: All settings configurable via environment
- **Command Line Flags**: Host and port via command line arguments
- **Default Values**: Sensible defaults for all settings

### Configuration Example

```go
package main

import (
    "context"
    "time"
    
    "github.com/aleal/zero/pkg/config"
    "github.com/aleal/zero/pkg/server"
)

func main() {
    ctx := context.Background()
    
    // Create custom configuration
    cfg := config.Default()
    cfg.SetHost("0.0.0.0")
    cfg.SetPort(3000)
    cfg.SetReadTimeout(10 * time.Second)
    cfg.SetWriteTimeout(30 * time.Second)
    cfg.SetAllowedOrigins([]string{"https://example.com"})
    
    // Create server with custom config
    server := server.NewServerWithConfig(ctx, cfg)
    
    // Add routes and start
    server.Start()
}
```

## 🛡️ Security Features

- **Input Validation**: All request parameters are sanitized and validated
- **CORS Support**: Configurable Cross-Origin Resource Sharing
- **Request Logging**: Comprehensive request/response logging with context
- **Panic Recovery**: Automatic recovery from panics with proper error responses
- **Type Safety**: Strong typing prevents common security issues
- **Context Isolation**: Request-specific context prevents data leakage

## 📊 Performance

- **Low Latency**: Optimized for minimal response times
- **High Throughput**: Efficient handling of concurrent requests
- **Memory Efficient**: Minimal memory allocation patterns
- **CPU Optimized**: Efficient goroutine usage and context switching

## 🔍 Built-in Endpoints

The example server includes several demonstration endpoints:

- `GET /health` - Health check endpoint
- `GET /hello` - Hello world with timestamp
- `GET /users` - List of example users
- `GET /users/{id}` - Get user by ID (path parameter example)
- `POST /users` - Create new user (JSON body parsing example)
- `GET /status` - Server status information
- `GET /static` - Static file serving example
- `GET /` - HTML documentation page

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Check for common issues
go vet ./...
```

## 🚀 Deployment

### Using in Your Application

```go
package main

import (
    "context"
    "github.com/aleal/zero/pkg/server"
    "github.com/aleal/zero/pkg/request"
    "github.com/aleal/zero/pkg/response"
)

func main() {
    ctx := context.Background()
    server := server.NewServer(ctx)
    
    // Add your application routes
    server.Get("/api/users", handleGetUsers)
    server.Post("/api/users", handleCreateUser)
    
    // Start the server
    server.Start()
}

func handleGetUsers(rctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Your handler logic here
    response.WriteJSON(w, http.StatusOK, []map[string]interface{}{
        {"id": 1, "name": "Alice"},
        {"id": 2, "name": "Bob"},
    })
}
```

### Docker Example

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server examples/server/cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8000
CMD ["./server"]
```

### Systemd Service

```ini
[Unit]
Description=Your App with Zero HTTP Server
After=network.target

[Service]
Type=simple
User=yourapp
WorkingDirectory=/opt/yourapp
ExecStart=/opt/yourapp/yourapp
Restart=always
RestartSec=5
Environment=ZERO_HOST=0.0.0.0
Environment=ZERO_PORT=8000

[Install]
WantedBy=multi-user.target
```

## 📚 API Documentation

### Server Interface

```go
type Server interface {
    Get(pattern string, handler request.Handler)
    Post(pattern string, handler request.Handler)
    Put(pattern string, handler request.Handler)
    Delete(pattern string, handler request.Handler)
    Patch(pattern string, handler request.Handler)
    Handle(pattern string, method string, handler request.Handler)
    Middlewares(middlewares ...middlewares.Middleware)
    Start()
}
```

### Request Utilities

```go
    // Parse JSON request body
    parser.ParseJSONBody(r, &user)

    // Get path parameters
    id := request.GetPathParam(r, "id")

    // Get query parameters
    page := request.GetQueryParam(r, "page")

    // Type-safe parameter parsing
    var limit int
    parser.ParseString("100", &limit)
```

### Response Utilities

```go
// Write JSON response
response.WriteJSON(w, http.StatusOK, data)

// Write error response
response.WriteError(w, http.StatusBadRequest, err)

// Set custom headers
response.SetHeader(w, "X-Custom-Header", "value")
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards and conventions
- Add comprehensive comments for all exported functions and types
- Include tests for new functionality
- Update documentation as needed
- Ensure all code passes `go vet` and `go fmt`

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with Go's excellent standard library
- Inspired by the simplicity and performance of minimal HTTP servers
- Thanks to the Go community for best practices and patterns
- Comprehensive documentation following Go standards

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/aleal/zero/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aleal/zero/discussions)
- **Documentation**: [GoDoc](https://pkg.go.dev/github.com/aleal/zero)

---

**Zero** - Because sometimes less is more. 🚀 