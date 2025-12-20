# Middleware Documentation

## HttpDebug Middleware

### Description
The `HttpDebug` middleware logs details about incoming HTTP requests and their responses. It is useful for debugging and monitoring request handling in the application.

### Features
- Logs the route name, HTTP method, path, response status code, and the time taken to process the request.

### Usage
To use the `HttpDebug` middleware, include it in your middleware chain:

```go
import "github.com/foxie-io/ng/example/basic/middlewares"

// Initialize the middleware
httpDebug := middlewares.HttpDebug{}

// Use it in your application
app.Use(httpDebug.Use)
```

### Example Log Output
```
GET /api/v1/resource 200 GET /api/v1/resource 15ms
```

---