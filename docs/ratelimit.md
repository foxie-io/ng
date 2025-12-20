# Middleware Documentation

## RateLimiter Middleware

### Description
The `RateLimiter` middleware enforces rate limiting on incoming HTTP requests. It tracks client requests and ensures that they do not exceed a predefined limit within a specified time window.

### Features
- Configurable request limits and time windows.
- Customizable client identification and error handling.
- Periodic cleanup of expired client data to prevent memory leaks.

### Usage
To use the `RateLimiter` middleware, initialize it with a configuration and include it in your middleware chain:

```go
import (
    "github.com/foxie-io/ng/example/basic/middlewares/limiter"
    "time"
)

// Create a rate limiter configuration
config := &limiter.Config{
    Limit:      100, // Maximum requests per window
    Window:     time.Minute, // Time window duration
    GenerateID: func(ctx context.Context) string { return "client-id" },
    ErrorHandler: func(ctx context.Context) error {
        return nghttp.NewErrTooManyRequests()
    },
}

// Initialize the middleware
rateLimiter := limiter.NewRateLimiter(config)

// Use it in your application
app.Use(rateLimiter.Use)
```

### Example
The middleware can be used to enforce rate limits on specific routes or globally across the application. For example:
```
POST /api/v1/resource
{
  "error": "Too Many Requests"
}
```

---