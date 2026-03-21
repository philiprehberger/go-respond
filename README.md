# go-respond

[![CI](https://github.com/philiprehberger/go-respond/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-respond/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-respond.svg)](https://pkg.go.dev/github.com/philiprehberger/go-respond) [![License](https://img.shields.io/github/license/philiprehberger/go-respond)](LICENSE)

HTTP JSON response helpers for Go. Write JSON responses in one line instead of five

## Installation

```bash
go get github.com/philiprehberger/go-respond
```

## Usage

```go
package main

import (
    "net/http"

    "github.com/philiprehberger/go-respond"
)

func main() {
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        user := map[string]string{"id": "1", "name": "Alice"}
        respond.OK(w, user)
    })

    http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
        item := map[string]string{"id": "42", "title": "New Item"}
        respond.Created(w, item)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        respond.JSON(w, http.StatusOK, map[string]string{"status": "healthy"})
    })

    http.ListenAndServe(":8080", nil)
}
```

### Error Responses

```go
// Simple error
respond.Error(w, http.StatusNotFound, "resource not found")
// {"error":{"status":404,"message":"resource not found"}}

// Error with details
respond.ErrorWithDetails(w, http.StatusBadRequest, "validation failed", map[string]string{
    "field":  "email",
    "reason": "invalid format",
})
// {"error":{"status":400,"message":"validation failed","details":{"field":"email","reason":"invalid format"}}}
```

### Validation Errors

```go
import respond "github.com/philiprehberger/go-respond"

respond.ValidationError(w, map[string]string{
    "email": "is required",
    "age":   "must be positive",
})
// {"error":"Validation failed","details":{"email":"is required","age":"must be positive"}}
```

### Paginated Responses

```go
import respond "github.com/philiprehberger/go-respond"

respond.Paginated(w, users, 100, 1, 20)
// {"data":[...],"meta":{"total":100,"page":1,"pageSize":20,"pages":5}}
```

### Problem Details (RFC 9457)

```go
respond.Problem(w, http.StatusForbidden,
    respond.WithType("https://example.com/problems/forbidden"),
    respond.WithTitle("Forbidden"),
    respond.WithDetail("You do not have access to this resource"),
    respond.WithInstance("/accounts/12345"),
    respond.WithExtension("account_id", "12345"),
)
// Content-Type: application/problem+json
// {"type":"https://example.com/problems/forbidden","title":"Forbidden","status":403,"detail":"You do not have access to this resource","instance":"/accounts/12345","account_id":"12345"}
```

## API

| Function | Description |
|----------|-------------|
| `JSON(w, status, data)` | Write JSON response with status code |
| `OK(w, data)` | Write 200 JSON response |
| `Created(w, data)` | Write 201 JSON response |
| `NoContent(w)` | Write 204 response |
| `Error(w, status, message)` | Write structured error response |
| `ErrorWithDetails(w, status, message, details)` | Write error with details |
| `ValidationError(w, errors)` | Write 422 with field validation errors |
| `Paginated[T](w, items, total, page, pageSize)` | Write 200 with pagination metadata |
| `Accepted(w, data)` | Write 202 Accepted response |
| `Problem(w, status, opts...)` | Write RFC 9457 Problem Details |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
