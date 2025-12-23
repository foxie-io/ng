# Advanced Example for NG Framework

This directory contains an advanced example of using the NG framework. It demonstrates how to structure a complex application with multiple components, middleware, guards, interceptors, and adapters.

## Project Structure

```
example/advance/
├── adapter/                     # Contains adapters for different frameworks and utilities
│   ├── echo.go                  # Adapter for the Echo framework
│   ├── stats.go                 # Middleware for collecting statistics
│   └── reqs/                    # Request-related utilities
│       ├── binder.go            # Custom request binder
│       └── validator.go         # Custom request validator
├── components/                  # Application components (modules)
│   ├── orders/                  # Orders module
│   │   ├── order.controller.go  # Controller for order-related routes
│   │   ├── order.module.go      # Module definition for orders
│   │   ├── order.service.go     # Business logic for orders
│   │   └── dtos/                # Data Transfer Objects for orders
│   │       ├── order.request.go # DTO for order requests
│   │       └── order.response.go# DTO for order responses
│   └── users/                   # Users module
│       ├── user.controller.go   # Controller for user-related routes
│       ├── user.module.go       # Module definition for users
│       ├── user.service.go      # Business logic for users
│       └── dtos/                # Data Transfer Objects for users
│           ├── user.request.go  # DTO for user requests
│           └── user.response.go # DTO for user responses
├── dal/                         # Data Access Layer
│   ├── order.dao.go             # DAO for orders
│   └── user.dao.go              # DAO for users
├── docs/                        # Documentation files
│   ├── docs.go                  # Swagger initialization
│   ├── swagger.json             # Swagger JSON file
│   └── swagger.yaml             # Swagger YAML file
├── models/                      # Database models
│   ├── order.model.go           # Model for orders
│   └── user.model.go            # Model for users
├── router/                      # Routing configuration
│   ├── grouper.go               # Route grouping logic
│   ├── router.go                # Main router setup
│   └── swagger.go               # Swagger route setup
├── go.mod                       # Go module file
├── go.sum                       # Go dependencies checksum file
├── main.go                      # Application entry point
└── Makefile                     # Build and run automation
```

## Dependencies

This example requires the following dependencies:

- `github.com/foxie-io/ng v0.3.0`: The core NG framework for building modular Go applications.
- `github.com/foxie-io/gormqs`: A library for simple query building and execution with GORM.
- `go.uber.org/fx v1.24.0`: A dependency injection framework for Go.
- `gorm.io/gorm v1.31.1`: An ORM library for Go.
- `github.com/labstack/echo/v4 v4.14.0`: A high-performance, extensible, and minimalist Go web framework.
- `github.com/swaggo/swag v1.16.6`: A library for generating Swagger documentation for Go applications.
- `github.com/MarceloPetrucio/go-scalar-api-reference`: API client UI from generated Swagger documentation.

## Features Demonstrated

1. **Controllers**:

   - Modular design for managing routes and handlers.
   - Example controllers for `orders` and `users`.

2. **Middleware**:

   - Custom middleware for logging, validation, and statistics.

3. **Guards**:

   - Guards for enforcing rules like rate limiting and authentication.

4. **Interceptors**:

   - Interceptors for transforming requests and responses.

5. **Adapters**:

   - Dynamic adapters for `echo`, `fiber`, and `http.ServeMux`.

6. **Data Access Layer (DAL)**:

   - Example DAOs for `orders` and `users`.

7. **Models and DTOs**:

   - Models for database entities.
   - DTOs for request and response validation.

8. **Swagger Documentation**:
   - Swagger files for API documentation.

## Running the Example

1. Install dependencies:

   ```bash
   go mod tidy
   ```

2. Run the application:

   ```bash
   go run main.go
   ```

3. Access the application:
   - API endpoints: `http://localhost:8080`
   - Swagger documentation: `http://localhost:8080/swagger`

## Learn More

Refer to the [NG Framework Documentation](../../README.md) for detailed information on features and usage.
