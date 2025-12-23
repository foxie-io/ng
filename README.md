# NestGo (NG) Framework

<div align="center">

A lightweight and modular Go framework inspired by NestJS, designed to simplify the development of scalable and maintainable web applications.

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/foxie-io/ng)](https://goreportcard.com/report/github.com/foxie-io/ng)

[Features](#features) •
[Installation](#installation) •
[Quick Start](#quick-start) •
[Examples](#examples) •
[Documentation](#core-concepts) •
[License](#license)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
  - [Request Flow](#request-flow)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Basic Application](#basic-application)
  - [Using Adapters](#using-adapters)
- [Examples](#examples)
- [Core Concepts](#core-concepts)
  - [Controllers](#controllers)
  - [Middleware](#middleware)
  - [Guards](#guards)
  - [Interceptors](#interceptors)
  - [Metadata](#metadata)
  - [Context Management](#context-management)
  - [Skippers](#skippers)
- [Advanced Topics](#advanced-topics)
  - [Sub-Applications](#sub-applications)
  - [Custom Adapters](#custom-adapters)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

**NestGo (NG)** brings the familiar, modular architecture of NestJS to the Go ecosystem. It provides a structured approach to building web applications with support for:

- **Controllers** for organizing routes and business logic
- **Middleware** for cross-cutting concerns like logging and authentication
- **Guards** for authorization and access control
- **Interceptors** for request/response transformation
- **Dynamic Adapters** for seamless integration with popular Go HTTP frameworks

Whether you're building a simple REST API or a complex microservice, NG helps you write clean, maintainable, and testable code.

---

## Features

- ✨ **NestJS-Inspired Architecture** - Familiar patterns for TypeScript developers transitioning to Go
- 🔌 **Dynamic HTTP Adapters** - Native support for Echo, Fiber, Gin, Chi, and standard `http.ServeMux`
- 🛡️ **Guards & Interceptors** - Built-in support for authentication, rate limiting, and request transformation
- 🎯 **Type-Safe Context Management** - Generic-based context storage with `ng.Store` and `ng.Load`
- 📦 **Modular Design** - Organize code into controllers, services, and modules
- 🔧 **Metadata System** - Attach configuration and behavior to routes dynamically
- **Comprehensive Examples** - Learn from basic to advanced use cases

---

## Architecture

### Request Flow

NestGo processes incoming HTTP requests through a well-defined pipeline:

```
┌─────────────────┐
│ Incoming Request│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Middleware    │  ◄── Logging, CORS, Body Parsing
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Guards      │  ◄── Authentication, Rate Limiting, Authorization
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Interceptors   │  ◄── Pre-processing, Validation
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Handler     │  ◄── Business Logic
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Interceptors   │  ◄── Post-processing, Transformation
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Response     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     Client      │
└─────────────────┘
```

**Pipeline Stages:**

1. **Middleware** - Execute first; handle logging, CORS, request parsing
2. **Guards** - Enforce access control and rate limits
3. **Interceptors (Pre)** - Transform or validate requests
4. **Handler** - Execute core business logic
5. **Interceptors (Post)** - Transform responses
6. **Response** - Send final response to client

---

## Installation

Add NestGo to your project:

```bash
go get github.com/foxie-io/ng
```

**Requirements:**

- Go 1.18 or higher (for generics support)

---

## Quick Start

### Basic Application

Create a simple health check API in minutes:

```go
package main

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	ngadapter "github.com/foxie-io/ng/adapter"
)

// Define a controller
type HealthController struct {
	ng.DefaultControllerInitializer
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/health"),
	)
}

// GET /health
func (c *HealthController) Index() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.WithHandler(func(ctx context.Context) error {
			return ng.Respond(ctx, nghttp.NewResponse("I am healthy!"))
		}),
	)
}

// GET /health/db
func (c *HealthController) DBCheck() ng.Route {
	return ng.NewRoute(http.MethodGet, "/db",
		ng.WithHandler(func(ctx context.Context) error {
			// Simulate database check
			return ng.Respond(ctx, nghttp.NewResponse(map[string]string{
				"status": "connected",
				"db":     "postgres",
			}))
		}),
	)
}

func main() {
	// Create application
	app := ng.NewApp(
		ng.WithResponseHandler(ngadapter.ServeMuxResponseHandler),
	)

	// Register controllers
	app.AddController(NewHealthController())

	// Build the application
	app.Build()

	// Setup HTTP server
	mux := http.NewServeMux()
	ngadapter.ServeMuxRegisterRoutes(app, mux)

	// Start server
	http.ListenAndServe(":8080", mux)
}
```

**Test it:**

```bash
curl http://localhost:8080/health
# Response: "I am healthy!"

curl http://localhost:8080/health/db
# Response: {"status":"connected","db":"postgres"}
```

### Using Adapters

NestGo supports multiple HTTP frameworks through adapters. Here's an example using Echo:

```go
package main

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
)

func main() {
	// Create NG application
	app := ng.NewApp()
	app.AddController(NewHealthController())
	app.Build()

	// Use Echo adapter
	e := echo.New()

	// Register NG routes with Echo
	RegisterEchoRoutes(app, e)

	// Start Echo server
	e.Start(":8080")
}
```

**Supported Adapters:**

- Standard `http.ServeMux`
- [Echo](https://echo.labstack.com/)
- [Fiber](https://gofiber.io/)
- [Gin](https://gin-gonic.com/)
- [Chi](https://go-chi.io/)

> **Note:** Typically, you'll use one adapter per application. Multiple adapters are useful for migration scenarios or serving the same routes on different ports.

---

## Examples

Explore real-world examples in the [`examples`](https://github.com/foxie-io/nestgo/tree/main/example) directory:

| Example                                                                       | Description                 | Key Features                               |
| ----------------------------------------------------------------------------- | --------------------------- | ------------------------------------------ |
| **[basic](https://github.com/foxie-io/nestgo/tree/main/example/basic)**       | Simple CRUD application     | Controllers, Middleware, Multiple adapters |
| **[advanced](https://github.com/foxie-io/nestgo/tree/main/example/advanced)** | Production-ready structure  | Guards, Interceptors, DAL, DTOs, Swagger   |
| **[chi](https://github.com/foxie-io/nestgo/tree/main/example/chi)**           | Chi router integration      | Chi adapter usage                          |
| **[echo](https://github.com/foxie-io/nestgo/tree/main/example/echo)**         | Echo framework integration  | Echo adapter, middleware                   |
| **[fiber](https://github.com/foxie-io/nestgo/tree/main/example/fiber)**       | Fiber framework integration | Fiber adapter, high performance            |
| **[gin](https://github.com/foxie-io/nestgo/tree/main/example/gin)**           | Gin framework integration   | Gin adapter, JSON handling                 |
| **[http](https://github.com/foxie-io/nestgo/tree/main/example/http)**         | Standard library only       | Native `http.ServeMux`, zero dependencies  |

Each example includes:

- Complete `main.go` with setup
- Controllers and route definitions
- Middleware and guard implementations
- Adapter configurations
- README with running instructions

---

## Core Concepts

### Controllers

Controllers organize related routes and encapsulate business logic. They provide a clean, modular structure for your application.

**Creating a Controller:**

```go
type UserController struct {
	ng.DefaultControllerInitializer
	userService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/users"),
		ng.WithMiddleware(LoggingMiddleware),
	)
}

// GET /users
func (c *UserController) List() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.WithHandler(func(ctx context.Context) error {
			users, err := c.userService.FindAll()
			if err != nil {
				return err
			}
			return ng.Respond(ctx, nghttp.NewResponse(users))
		}),
	)
}

// GET /users/:id
func (c *UserController) Get() ng.Route {
	return ng.NewRoute(http.MethodGet, "/:id",
		ng.WithHandler(func(ctx context.Context) error {
			id := ng.Param(ctx, "id")
			user, err := c.userService.FindByID(id)
			if err != nil {
				return err
			}
			return ng.Respond(ctx, nghttp.NewResponse(user))
		}),
	)
}

// POST /users
func (c *UserController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		ng.WithHandler(func(ctx context.Context) error {
			var input CreateUserDTO
			if err := ng.Bind(ctx, &input); err != nil {
				return err
			}

			user, err := c.userService.Create(input)
			if err != nil {
				return err
			}
			return ng.Respond(ctx, nghttp.NewResponse(user))
		}),
	)
}
```

**Key Points:**

- Controllers group related routes under a common prefix
- Use dependency injection for services
- Each method returns a `ng.Route`
- Routes can have their own middleware, guards, and interceptors

---

### Middleware

Middleware functions execute before guards and interceptors. They're perfect for logging, CORS, authentication, and request parsing.

**Creating Middleware:**

```go
type LoggingMiddleware struct {
	ng.DefaultID[LoggingMiddleware]
}

func (m LoggingMiddleware) Use(ctx context.Context, next ng.Handler) {
	start := time.Now()

	// Log incoming request
	log.Printf("[%s] %s", ng.Method(ctx), ng.Path(ctx))

    defer func() {
        // Log response time
        log.Printf("Completed in %v", time.Since(start))
    }()

	// Call next handler
    next(ctx)

    // no return because it is guard job to brock request
    // but can be force to stop here by ng.ThrowResponse or ng.ThrowAny
}

// Apply to entire application
app := ng.NewApp(
	ng.WithMiddleware(LoggingMiddleware{}),
)

// Apply to specific controller
func (c *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/users"),
		ng.WithMiddleware(LoggingMiddleware{}),
	)
}

// Apply to specific route
func (c *UserController) RouteWithLog() ng.Route {
	return ng.NewRoute(http.MethodGet, "/route-with-log",
		ng.WithHandler(handler), // controller middleware also applied
	)
}

func (c *UserController) RouteWithoutLog() ng.Route {
	return ng.NewRoute(http.MethodGet, "/route-without-log",
        ng.Skip(LoggingMiddleware{}), // skip this middleware
		ng.WithHandler(handler),
	)
}
```

---

### Guards

Guards determine whether a request should proceed. They're executed after middleware and are ideal for authentication, authorization, and rate limiting.

**Creating a Guard:**

```go
type AuthGuard struct {
	ng.DefaultID[AuthGuard]
}

func (g AuthGuard) Allow(ctx context.Context) {
	// Extract token from request
	token := ng.Header(ctx, "Authorization")
	user, err := validateToken(token)
	if err != nil {
		return nghttp.NewError(http.StatusUnauthorized, "Invalid token")
	}

	// Store user in context for handlers
	ng.Store(ctx, user) // type: (*User)

    // can be extract by ng.Load[*User](ctx) in handler
}

// Apply guard
app := ng.NewApp(
	ng.WithGuards(AuthGuard{}),
)

func (c *UserController) Login() ng.Route {
	return ng.NewRoute(http.MethodGet, "/login",
        ng.SkipAllGuards(), // become public route
        ng.Skip(AuthGuard{}),// or skip only auth guard
		ng.WithHandler(handler),
	)
}
```

**Rate Limiting Guard:**

```go
type RateLimitGuard struct {
	ng.DefaultID[RateLimitGuard]
}


type LimitConfig struct {
    Limit  int
    Window time.Duration
}

// use type as key is safer than string
type limitConfigKey struct {
    ng.TypeKey[limitConfigKey]
}

func (g *RateLimitGuard) Allow(ctx context.Context) error {
    config,_ := ng.GetContext(ctx).Route().Core().Metadata(limitConfigKey{}).(*LimitConfig)
    if config == nil {
        config = defaultConfig
    }

    // check rate limit
	return applyRateLimit(ctx, config)
}

// declare your own metadata helper
func WithLimitConfig(limit int, window time.Duration) ng.Option {
    return ng.WithMetadata(limitConfigKey{}, &LimitConfig{
        Limit:  limit,
        Window: window,
    })
}


// Apply rete limit guard
app := ng.NewApp(
    // also can use ng.Options to merge options
    // UseRateLimit := ng.Options(ng.WithGuard(...),ng.WithLimitConfig(...))
	ng.WithGuards(RateLimitGuard{}),
    ng.WithLimitConfig(100, time.Minute)
    ...
)

// override rate limit config per route or controller
func (c *UserController) Create() ng.Route {
    return ng.NewRoute(http.MethodPost, "/",
        ng.WithMetadata("RateLimitKey", &LimitConfig{
            Limit:  5,
            Window: time.Minute,
        }),
        ...
    )
}
```

---

### Interceptors

Interceptors transform requests before they reach handlers and responses after handlers execute. They're perfect for validation, transformation, and logging.

**Creating an Interceptor:**

```go
type TransformInterceptor struct {
	ng.DefaultID[TransformInterceptor]
}

func (i TransformInterceptor) Intercept(ctx context.Context, next ng.Handler)  {
	// Pre-processing
    start := time.Now()

    defer func() {
        rctx := ng.GetContext(ctx)
        // Post-processing: log response time

        transformedResp := TransformResponse(rctx.Response())

        rctx.SetResponse(transformedResp)

        logResponseTime(start)
    }()

    next(ctx)
}

// Apply interceptor
func (c *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/users"),
		ng.WithInterceptors(TransformInterceptor{}),
	)
}
```

**Response Wrapper Interceptor:**

```go
type ResponseWrapperInterceptor struct {
	ng.DefaultID[ResponseWrapperInterceptor]
}
func (i ResponseWrapperInterceptor) Intercept(ctx context.Context, next ng.Handler) {
    defer func() {
        rctx := ng.GetContext(ctx)
        response := rctx.Response()

        if response.StatusCode() >= 400 {
            // Do not wrap error responses
            return
        }

        wrappedResponse := nghttp.NewResponse(map[string]interface{}{
            "data": response.Response(),
            "meta": map[string]interface{}{
                "status":  response.StatusCode(),
                "message": "success",
            },
        })

        rctx.SetResponse(wrappedResponse)
    }()

    next(ctx)
}
```

---

### Metadata

Metadata allows dynamic configuration of routes, controllers, and components. It's inspired by NestJS decorators.

**Using Metadata:**

```go
// Define metadata keys
type rateLimitKey struct {
    ng.TypeKey[rateLimitKey]
}

// Attach metadata to routes
func (c *UserController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		ng.WithMetadata(rateLimitKey{}, &RateLimitConfig{
			Limit:  10,
			Window: time.Minute,
		}),
		ng.WithHandler(handler),
	)
}

// Access metadata in guards/middleware
func (g RateLimitGuard) Allow(ctx context.Context) error {
	config, exists := ng.GetContext(ctx).Route().Core().Metadata(rateLimitKey{}).(*RateLimitConfig)
	if !exists {
		// Use default config
		config = defaultConfig
	}

	// Apply rate limiting based on config
	return applyRateLimit(ctx, config)
}
```

**Metadata Hierarchy:**

Metadata can be set at three levels:

1. **Application Level** - Applies to all routes
2. **Controller Level** - Applies to all routes in the controller
3. **Route Level** - Applies to specific route

Route-level metadata overrides controller-level, which overrides application-level.

**Helper Functions:**

```go
// Custom metadata helper
func WithRateLimit(limit int, window time.Duration) ng.Option {
	return ng.WithMetadata(RateLimitKey, &RateLimitConfig{
		Limit:  limit,
		Window: window,
	})
}

// Usage
func (c *UserController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		WithRateLimit(10, time.Minute),
		ng.WithHandler(handler),
	)
}
```

---

### Context Management

NestGo provides type-safe context utilities for storing and retrieving request-scoped data.

**Storing Data:**

```go
type User struct {
	ID   int
	Name string
	Role string
}

// In middleware or guard
func AuthMiddleware(ctx context.Context, next ng.Handler) error {
	user := User{ID: 1, Name: "John Doe", Role: "admin"}
	ng.Store(ctx, user)
	return next(ctx)
}
```

**Loading Data:**

```go
// In handler
func UserHandler(ctx context.Context) error {
	// Load with check
	user, exists := ng.Load[User](ctx)
	if !exists {
		return errors.New("user not found in context")
	}

	return ng.Respond(ctx, nghttp.NewResponse(user))
}

// Must load (panics if not exists)
func AdminHandler(ctx context.Context) error {
	user := ng.MustLoad[User](ctx)
	// ...
}

// Load with default
func DefaultHandler(ctx context.Context) error {
	user,loaded := ng.LoadOrStore(ctx, User{Name: "Guest", Role: "guest"})
	// ...
}
```

**Available Functions:**

```go
// Store value in context
ng.Store[T](ctx context.Context, value T)

// Load value from context
ng.Load[T](ctx context.Context) (value T, exists bool)

// Load value or panic if not exists
ng.MustLoad[T](ctx context.Context) T

// Load value or store default
ng.LoadOrStore[T](ctx context.Context, defaultValue T) T
```

---

### Skippers

Skippers allow you to conditionally bypass middleware, guards, or interceptors for specific routes.

**Using DefaultID:**

```go
type AuthGuard struct {
	ng.DefaultID[AuthGuard]
}

// Skip auth guard for public routes
func (c *UserController) PublicProfile() ng.Route {
	return ng.NewRoute(http.MethodGet, "/public/:id",
		ng.WithSkip(AuthGuard{}),
		ng.WithHandler(handler),
	)
}
```

**Skip All Guards:**

```go
func (c *HealthController) Check() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.SkipAllGuards(),
		ng.WithHandler(handler),
	)
}
```

**Custom ID Implementation:**

```go
type CustomMiddleware struct{}

func (m CustomMiddleware) NgID() string {
	return "CustomMiddleware"
}

// Skip by ID
ng.WithSkip(CustomMiddleware{})
```

**Key Points:**

- Use `DefaultID[T]` for automatic ID generation
- Use `WithSkip()` to skip specific components
- Use `SkipAllGuards()` for public endpoints
- Skippers work with middleware, guards, and interceptors

---

## Advanced Topics

### Sub-Applications

Organize large projects by mounting sub-applications:

```go
// Create sub-apps for different modules
adminApp := ng.NewApp(
	ng.WithPrefix("/admin"),
	ng.WithGuards(AdminGuard{}),
)
adminApp.AddController(NewAdminController())

apiApp := ng.NewApp(
	ng.WithPrefix("/api"),
	ng.WithMiddleware(RateLimitMiddleware{}),
)
apiApp.AddController(NewAPIController())

// Main application
app := ng.NewApp()
app.AddSubApp(adminApp)
app.AddSubApp(apiApp)
app.Build()
```

### Custom Adapters

Create adapters for other HTTP frameworks:

```go
func CustomRegisterRoutes(app *ng.App, router *CustomRouter) {
	for _, route := range app.Routes() {
		router.Handle(route.Method, route.Path, tranformHandler(route.handler))
	}
}
```

---

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

**Development:**

```bash
# Clone the repository
git clone https://github.com/foxie-io/ng.git

# Run tests
go test ./...

# Run examples
cd example/basic
go run main.go
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

<div align="center">

**[⬆ Back to Top](#nestgo-ng-framework)**

</div>
