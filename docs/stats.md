# Middleware Documentation

## Stats Middleware

### Description
The `Stats` middleware collects and provides statistics about incoming HTTP requests. It is useful for monitoring and analyzing application performance.

### Features
- Tracks request counts, response times, and other metrics.
- Provides an interface to retrieve collected statistics.

### Usage
To use the `Stats` middleware, include it in your middleware chain:

```go
import "github.com/foxie-io/ng/example/basic/middlewares"

// Initialize the middleware
stats := middlewares.Stats{}

// Use it in your application
app.Use(stats.Use)
```

### Example
The middleware can be configured to expose an endpoint for retrieving statistics, such as:
```
GET /stats
{
  "requests": 100,
  "average_response_time": "20ms"
}
```

---