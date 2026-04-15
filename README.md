# Zero - Pure Golang HTTP Server Library

A simple, lean, and blazingly fast HTTP server library built with pure Go using only native libraries. Designed for maximum performance, security, and minimal resource footprint. Zero provides a lightweight foundation for building high-performance HTTP servers with comprehensive middleware support and graceful shutdown capabilities.

## Features

- **Zero Dependencies**: Uses only Go's standard library - no external packages
- **Lightning Fast**: Optimized for high-performance request handling
- **Memory Efficient**: Minimal memory footprint and garbage collection pressure
- **Security First**: Built with security best practices and input validation
- **Simple & Clean**: Easy to understand, maintain, and extend
- **Production Ready**: Robust error handling, structured logging, and graceful shutdown
- **Middleware Support**: Built-in CORS, `log/slog` request logging, panic recovery, and request ID generation
- **Type-Safe**: Generic-based type-safe parameter parsing
- **Context Integration**: Full support for request context, structured logging, and request ID tracing

## Requirements

- Go 1.24.6 or higher
- No external dependencies required

## Installation

```bash
go get github.com/aleal/zero
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "net/http"

    "github.com/aleal/zero/pkg/response"
    zero "github.com/aleal/zero/pkg/server"
)

func main() {
    ctx := context.Background()

    srv := zero.New(ctx, zero.WithDefaultMiddlewares())

    srv.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
        response.WriteJSON(w, http.StatusOK, map[string]string{
            "message": "Hello from Zero!",
        })
    })

    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

`New` registers `GET /health` automatically. `Start` blocks until a signal is received, then shuts down gracefully and returns any error.

## Environment Variables

All configuration is done via environment variables and programmatic options. There are no command-line flags — the library never calls `flag.Parse()`.

| Variable | Default | Description |
|----------|---------|-------------|
| `ZERO_HOST` | `localhost` | Server host |
| `ZERO_PORT` | `8000` | Server port |
| `ZERO_READ_TIMEOUT` | `5s` | HTTP read timeout |
| `ZERO_WRITE_TIMEOUT` | `15s` | HTTP write timeout |
| `ZERO_IDLE_TIMEOUT` | `60s` | HTTP idle timeout |
| `ZERO_MAX_JSON_REQUEST_BODY_SIZE` | `1048576` | Max JSON body size in bytes (1 MB) |
| `ZERO_MAX_UPLOADED_FILE_SIZE` | `10485760` | Max uploaded file size in bytes (10 MB) |
| `ZERO_LOG_LEVEL` | `INFO` | Log level: `DEBUG`, `INFO`, `WARNING`, `ERROR` |

Invalid env values are logged at `WARN` level and the default is used — the server always starts. Middleware (CORS, logging, recovery) is configured in code via server options, not environment variables.

## Architecture

```
zero/
├── pkg/
│   ├── server/            # Core server, interface, and options
│   ├── config/            # Configuration (env vars + defaults)
│   ├── log/               # slog-based structured logging
│   ├── metadata/           # Library version (auto-detected from Go module info)
│   ├── middlewares/        # CORS, logging, recovery, middleware chain
│   ├── requestid/          # Hostname-counter request ID generation
│   ├── request/            # Request utilities and builder
│   ├── response/           # Response utilities
│   └── parser/             # Type-safe parsing and JSON body decoding
├── internal/
│   └── handlers/           # Built-in handlers (health check)
├── examples/
│   └── server/cmd/         # Example server
└── README.md
```

## Configuration

### Programmatic

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/aleal/zero/pkg/config"
    "github.com/aleal/zero/pkg/middlewares"
    "github.com/aleal/zero/pkg/server"
)

func main() {
    ctx := context.Background()

    cfg := config.Load()
    cfg.SetHost("0.0.0.0")
    cfg.SetPort(3000)
    cfg.SetReadTimeout(10 * time.Second)
    cfg.SetWriteTimeout(30 * time.Second)
    cfg.SetMaxJSONBodySize(5 << 20) // 5 MB

    srv := server.New(ctx,
        server.WithConfig(cfg),
        server.WithCORS([]string{"https://example.com"}, middlewares.MiddlewarePriorityLow),
        server.WithDefaultLogging(),
        server.WithDefaultRecovery(),
    )

    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}
```

### Options

Common options passed to `server.New(ctx, ...)`:

| Option | Description |
|--------|-------------|
| `WithConfig(cfg)` | Use a custom `*config.Config` |
| `WithHost(host)` | Override host |
| `WithPort(port)` | Override port |
| `WithReadTimeout(d)` | Override read timeout |
| `WithWriteTimeout(d)` | Override write timeout |
| `WithIdleTimeout(d)` | Override idle timeout |
| `WithMaxUploadedFileSize(n)` | Override max uploaded file size |
| `WithDefaultMiddlewares()` | Apply default logging + CORS + recovery |
| `WithDefaultLogging()` | Structured JSON logging via `log/slog` |
| `WithDefaultCORS()` | CORS with wildcard origin |
| `WithDefaultRecovery()` | Panic recovery middleware |
| `WithCORS(origins, priority)` | CORS with specific origins |
| `WithLogging(logger, priority)` | Custom slog logger |
| `WithRecovery(priority)` | Recovery middleware |
| `WithMiddleware(mw, priority)` | Custom middleware |

## Server Interface

```go
type Zero interface {
    Get(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
    Post(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
    Put(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
    Delete(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
    Patch(pattern string, handler request.Handler, middlewares ...middlewares.Middleware)
    Handle(pattern string, method string, handler request.Handler, middlewares ...middlewares.Middleware)
    Handler() http.Handler
    Start() error
    Shutdown(ctx context.Context) error
}
```

- `Handler()` returns the underlying `http.Handler` — useful for `httptest.NewServer(srv.Handler())` in tests.
- `Start()` blocks until interrupted, then performs graceful shutdown and returns any error.
- `Shutdown(ctx)` triggers graceful shutdown programmatically.
- Route methods accept optional per-route middlewares.

## Request Utilities

```go
// Parse JSON body (respects ZERO_MAX_JSON_REQUEST_BODY_SIZE, or pass explicit limit)
if err := parser.ParseJSONBody(r.Body, &user); err != nil {
    // handle error
}

// With explicit max size
if err := parser.ParseJSONBody(r.Body, &user, cfg.MaxJSONBodySize); err != nil {
    // handle error
}

// Path parameters (Go 1.22+ routing, e.g. pattern "/users/{id}")
id := request.GetPathParam(r, "id")

// Query parameters
page := request.GetQueryParam(r, "page")
limit := request.GetParsedQueryParamOrDefault[int](r, "limit", 20)

// Type-safe parsing from string
var port int
_ = parser.ParseString("8080", &port)
```

## Request Builder

Fluent builder for constructing outbound HTTP requests. Supports JSON bodies, multipart form uploads, and custom headers.

```go
// JSON POST
req, err := request.NewBuilder(ctx, http.MethodPost, "https://api.example.com/users").
    WithBodyJSON(map[string]string{"name": "Alice"}).
    WithHeader("Authorization", "Bearer tok").
    Build()
if err != nil {
    // handle error
}

client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
```

```go
// Multipart file upload
file := request.NewUploadedFile(data, "photo.jpg", "image/jpeg")
req, err := request.NewBuilder(ctx, http.MethodPost, "https://api.example.com/upload").
    WithFormFile("avatar", file).
    WithFormField("description", "Profile photo").
    Build()
```

### Builder Interface

```go
type Builder interface {
    WithContext(ctx context.Context) Builder
    WithMethod(method string) Builder
    WithURL(url string) Builder
    WithHeader(key, value string) Builder
    WithFormFile(key string, file *UploadedFile) Builder
    WithFormField(key, value string) Builder
    WithBodyJSON(body any) Builder
    Build() (*http.Request, error)
}
```

## Response Utilities

```go
response.WriteJSON(w, http.StatusOK, data)
response.WriteError(w, http.StatusBadRequest, err)
response.WriteErrorMsg(w, http.StatusNotFound, "not found")
response.SetHeader(w, "X-Custom-Header", "value")
```

## Request IDs

The logging middleware generates request IDs in the format `hostname-counter` using an atomic counter. These are:

- Injected into the structured logger as `requestId`
- Stored in request context (retrieve with `requestid.FromContext(r.Context())`)
- Unique per host, monotonically increasing, zero crypto overhead

## Security

- **Input Validation**: All request parameters sanitized (null bytes, control characters removed)
- **CORS**: Wildcard origin uses literal `*` without credentials; explicit origins set `Access-Control-Allow-Credentials: true` with `Vary: Origin`
- **Panic Recovery**: Generic error returned to client; panic details logged server-side only
- **Body Limits**: JSON body parsing limited by `ZERO_MAX_JSON_REQUEST_BODY_SIZE` (default 1 MB); file uploads rejected above `ZERO_MAX_UPLOADED_FILE_SIZE` (default 10 MB)
- **Type Safety**: Generic-based parsing prevents common injection issues

## Built-in Endpoints

Every server instance registers `GET /health` automatically, returning:

```json
{"service": "zero", "version": "v0.1.0", "uptime": "1h23m45s"}
```

The version is detected automatically from Go's module build info (`runtime/debug.ReadBuildInfo`). When imported as a dependency via `go get github.com/aleal/zero@v0.1.0`, the tagged version appears in the banner and health endpoint with no manual steps. During local development it shows `devel`.

The example server adds: `/hello`, `/users`, `/users/{id}`, `/status`, `/static`, `/`.

## Testing

```bash
# Run all tests with coverage
make test

# Open HTML coverage report
make coverage

# Run benchmarks
make bench

# Check for issues
go vet ./...
```

Use `Handler()` in tests to exercise the actual server routing:

```go
func TestAPI(t *testing.T) {
    srv := server.New(context.Background())
    srv.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
        response.WriteJSON(w, http.StatusOK, map[string]string{"pong": "ok"})
    })

    ts := httptest.NewServer(srv.Handler())
    defer ts.Close()

    resp, _ := http.Get(ts.URL + "/ping")
    // assert resp...
}
```

## Deployment

### Docker

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

### Systemd

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

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go coding standards and conventions
- Add comments for all exported functions and types
- Include tests for new functionality
- Ensure all code passes `go vet` and `go fmt`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Zero** - Because sometimes less is more.
