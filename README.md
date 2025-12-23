# (NG) NestGo Framework Documentation

## Overview

NestGo is a lightweight and modular Go framework inspired by NestJS. It is designed to simplify the development of scalable and maintainable web applications. NG provides a structured approach to building applications with support for controllers, middleware, guards, interceptors, and dynamic adapters for various HTTP frameworks.

---

## Features

- **Inspired by NestJS**: Brings the modular and organized structure of NestJS to Go.
- **Middleware Support**: Easily integrate middlewares for logging, rate limiting, and more.
- **Dynamic Adapters**: Support for multiple HTTP frameworks like `echo`, `fiber`, and `http.ServeMux`.
- **Guards**: Implement guards to enforce rules such as rate limiting.
- **Interceptors**: Modify or transform requests and responses before or after they are processed by the handler.
- **Controllers**: Modular controller design for managing routes and handlers.
- **SubApps**: Support for sub-applications to organize large projects.
- **Customizable**: Flexible configuration options for middlewares, guards, interceptors, and response handlers.

---

## Request Flow

The NG framework processes incoming requests through the following sequence:

1. **Middleware**:

   - Middleware functions are executed first. They can modify the request, add metadata, or perform logging and monitoring.

2. **Guards**:

   - Guards are executed next. They determine whether the request is allowed to proceed based on custom rules (e.g., rate limiting, authentication).

3. **Interceptors**:

   - Interceptors are executed before and after the handler. They can modify the request before it reaches the handler or transform the response after the handler has processed the request.

4. **Handler**:

   - The handler processes the request and generates a response. This is the core business logic of the application.

5. **Response**:
   - The response is sent back to the client. Middleware or interceptors can modify the response before it is finalized.

### Diagram

```
Incoming Request → Middleware → Guards → Interceptors → Handler → Interceptors → Response → Client
```

---

## Installation

To use NG in your project, add it to your `go.mod` file:

```bash
go get github.com/foxie-io/ng
```

---

## Quick Start

### 1. Create a New Application with a Controller

```go
type HealthController struct {
	ng.DefaultControllerInitializer
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (con *HealthController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/health"),
	)
}

// curl localhost:8080/health
func (con *HealthController) Index() ng.Route {
	return ng.NewRoute(http.MethodGet, "/", ng.WithHandler(func(ctx context.Context) error {
		return ng.Respond(ctx, nghttp.NewResponse("I am ok!"))
	}))
}

// curl localhost:8080/health/db-check
func (con *HealthController) DBCheck() ng.Route {
	return ng.NewRoute(http.MethodGet, "/db-check", ng.WithHandler(func(ctx context.Context) error {
		return ng.Respond(ctx, nghttp.NewResponse("db is ok"))
	}))
}

func runApp() {
	app := ng.NewApp(
		ng.WithResponseHandler(ngadapter.ServeMuxResponseHandler),
	)

	// Add the HealthController
	app.AddController(NewHealthController())

	// Build and start the application
	app.Build()

	// Start the HTTP server
	mux := http.NewServeMux()
	ngadapter.ServeMuxRegisterRoutes(app, mux)

	// Listen and serve
	http.ListenAndServe(":8080", mux)
}
```

### 2. Add Middleware

```go
// example/basic/middlewares/debug.go
app := ng.NewApp(
	ng.WithMiddleware(
		middlewares.HttpDebug{},
	),
)
```

### 3. Add Guards

```go
// example/basic/middlewares/limiter/limiter.go
app := ng.NewApp(
	ng.WithGuards(
		limiter.New(&limiter.Config{
			Limit:  5,
			Window: time.Minute,
			GenerateID: func(ctx context.Context) string {
				return "client-id"
			},
		}),
	),
)
```

### 4. Add Interceptors

```go
// no example yet
package main

import (
	"context"
	"github.com/foxie-io/ng"
	"log"
)

type RequestTransformer struct {}

func NewRequestTransformer() *RequestTransformer {
	return &RequestTransformer{}
}

func (rt *RequestTransformer) Intercept(ctx context.Context, next ng.Handler) {
	// Before operation: Transform the request
	log.Println("Transforming request...")

    defer func() {
        // After operation
        // can modify response here if needed
        rc := ng.GetContext(ctx)
        resp := rc.GetResponse()

        // do sth with resp
        log.Printf("Response status: %d", resp.StatusCode())

        rc.SetResponse(resp)
    }()

	// Call the next handler in the chain
	next(ctx)
}

func main() {
	app := ng.NewApp(
		ng.WithInterceptors(
			NewRequestTransformer(),
		),
	)

	// Build and start the application
	app.Build()
}
```

### 5. Use Adapters

this is an example of using multiple adapters in the same application
normally, you would only use one adapter per application

```go
// example/basic/main.go
muxHttp := http.NewServeMux()
echoHttp := echo.New()
fiberHttp := fiber.New()

// Register routes to adapters
adapters.ServeMuxRegisterRoutes(app, muxHttp)
adapters.EchoRegisterRoutes(app, echoHttp)
adapters.FiberRegisterRoutes(app, fiberHttp)

// Start servers
go fiberHttp.Listen(":8082")
go http.ListenAndServe(":8080", muxHttp)
echoHttp.Start(":8081")
```

---

## Examples Directory

The `examples` directory contains various subdirectories showcasing how to use the NG framework with different setups and configurations. Below is a list of available examples:

- **basic/**: A simple example demonstrating the core features of the NG framework, including middleware, controllers, and adapters.
- **advance/**: An advanced example with a more complex project structure, demonstrating controllers, middleware, guards, interceptors, adapters, data access layers, models, DTOs, and Swagger documentation.
- **chi/**: Example using the `chi` router as the HTTP adapter.
- **fiber/**: Example using the `fiber` framework as the HTTP adapter.
- **fx/**: Example integrating NG with the `fx` dependency injection framework.
- **gin/**: Example using the `gin` framework as the HTTP adapter.
- **http/**: Example using the standard `http.ServeMux` as the HTTP adapter.

Each example contains its own `main.go` file and supporting components, middleware, and adapters. Navigate to the respective directories to explore the code and learn how to integrate NG into your projects.

---

## Documentation

### Controllers

Controllers group related routes and provide a modular way to manage application logic. Use `ng.NewController` to create a controller.

### Middleware

Middlewares are reusable components that process requests and responses. Examples include logging, rate limiting, and statistics collection.

### Guards

Guards enforce rules such as rate limiting, role guards. They are executed after middleware and before interceptors.

### Interceptors

Interceptors allow you to modify requests before they reach the handler and responses after they are processed by the handler. They are useful for tasks like logging, transforming data, or adding headers.
it executed after guards and before the handler.

### Metadata

Metadata allows you to attach additional information to routes, controllers, and other components. This pattern is inspired by NestJS and provides a flexible way to modify behavior of middleware,guards, interceptors dynamically based on metadata.

---

# Metadata in NG Framework

## Overview

Metadata in the NG framework allows you to attach additional information to routes, controllers, and other components. This pattern is inspired by NestJS and provides a flexible way to modify behavior dynamically based on metadata.

---

## What is Metadata?

Metadata is structured data that can be associated with routes, controllers, or other components in the NG framework. It is used to:

- Pass configuration options.
- Enable or disable middleware or guards dynamically.
- Attach additional information for runtime behavior.

---

## Metadata in Routes

You can attach metadata to routes using the `ng.WithMetadata` option. This metadata can then be accessed by middleware, guards, or other components.

### Example: Attaching Metadata to a Route

```go
package main

import (
	"context"
	"github.com/foxie-io/ng"
	"github.com/foxie-io/ng/http"
)

func main() {
	app := ng.NewApp(
        ng.WithPrefix("/api/v1/resource"),
		ng.WithGuards(
            limiter.New(&limiter.Config{
                Limit:  10,
                Window: time.Minute,
                GenerateID: func(ctx context.Context) string {
                    return ctx.Value("client-id").(string)
                },
            }),
        ),
    )

	app.AddController(NewHealthController())
	app.Build()
}
```

In this example:

- Metadata is attached to the route with the key `rateLimit`.
- Middleware or guards can access this metadata to enforce rate limiting dynamically.

---

## Metadata in Controllers & Routes

Controllers can also have metadata that applies to all routes within the controller. This is useful for setting global options for a group of routes.

### Example: Attaching Metadata to a Controller and Route

```go
type OrdersController struct {
	ng.DefaultControllerInitializer
}

func NewOrdersController() *OrdersController {
	return &OrdersController{}
}

func (con *OrdersController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/orders"),
        // Override app level metadata apply to all routes in this controller
        // Depend on ng.WithMetadata
        // file: example/basic/middlewares/limiter/option.go
		limiter.WithConfig(&limiter.Config{
			Limit:  10,
			Window: time.Minute,
		}),
	)
}

func (con *OrdersController) ViewOrder() ng.Route {
    return ng.NewRoute(http.MethodGet, "/view",
        // limit base on controller metadata
        ng.WithHandler(func(ctx context.Context) error {
            return ng.Respond(ctx, http.NewResponse("Order viewed"))
        }),
    )
}

func (con *OrdersController) CreateOrder() ng.Route {
    return ng.NewRoute(http.MethodPost, "/create",
        // Override app level and controller-level metadata
        limiter.WithConfig(&limiter.Config{
            Limit:  3,
            Window: time.Minute,
        }),
        ng.WithHandler(func(ctx context.Context) error {
            return ng.Respond(ctx, http.NewResponse("Order created"))
        }),
    )
}
```

## Skipper and ID Interface

The NG framework provides a mechanism to conditionally skip middleware, guards, or interceptors during request execution using the `ID` interface and skipper utilities.

### ID Interface

The `ID` interface is used to uniquely identify components that can be skipped. It ensures that each component has a stable and unique identifier.

```go
type ID interface {
    NgID() string
}

type AuthGuard struct {}
func (a AuthGuard) NgID() string {
    return "AuthGuard"
}
ng.WithSkipper(AuthGuard{})

```

### DefaultID

The `DefaultID` generic type allows skipping components by their concrete type. For example:

```go
type AuthGuard struct {
    ng.DefaultID[AuthGuard]
}
ng.WithSkipper(AuthGuard{})
```

This will prevent the `AuthGuard` from executing for the configured route.

### Skipper Utilities

- **WithSkip**: Attaches skipper metadata to a route or handler. Components implementing the `ID` interface can be skipped if their `NgID` matches the provided skipper IDs.

  ```go
  func WithSkip(skippers ...ID) Option
  ```

- **SkipAllGuards**: Skips the execution of all guards for a route or handler. Useful for public endpoints like health checks.

  ```go
  func SkipAllGuards() Option
  ```

### Corrected Example

To use `DefaultID` with a concrete type like `AuthGuard`, ensure that the type is defined as follows:

```go
type AuthGuard struct {
    ng.DefaultID[AuthGuard]
}
```

Then, you can skip the `AuthGuard` middleware as shown:

```go
app := ng.NewApp(
    ng.WithMiddleware(
        AuthGuard{},
    ),
    ng.WithSkipper(
        AuthGuard{},
    ),
)
```

This ensures that the `AuthGuard` middleware is properly defined and can be skipped for the configured route.

## Summary

- **Route Metadata**: Attach metadata to individual routes for dynamic behavior.
- **Controller Metadata**: Apply metadata to all routes within a controller.
- **Accessing Metadata**: Use `ng.GetContext` to retrieve metadata at runtime.

By leveraging metadata, the NG framework provides a powerful and flexible way to customize behavior dynamically, similar to the NestJS pattern.

---

## Context Management with ng.Store and ng.Load

The NG framework provides utilities for managing context-specific data using `ng.Store` and `ng.Load`. These functions simplify the process of attaching and retrieving values from the request context.

### ng.Store

The `ng.Store` function is used to attach a value to the context. This is particularly useful in middleware or adapters where you need to pass data (e.g., user information, request metadata) to subsequent handlers.

#### Function Signature

```go
func Store[T any](ctx context.Context, value T)
```

#### Example Usage

```go
// Middleware example
func AttachUserMiddleware(ctx context.Context, next ng.Handler) {
    user := User{ID: 1, Name: "John Doe"}
    ng.Store(ctx, user)
    next(ctx)
}
```

### ng.Load

The `ng.Load` function is used to retrieve a value from the context that was previously stored using `ng.Store`.

#### Function Signature

```go
func Load[T any](ctx context.Context) (value T, exists bool)
```

#### Example Usage

```go
// Guard example
func AdminGuard(ctx context.Context) error {
    user, exists := ng.Load[User](ctx)
    if !exists {
        return errors.New("user not found")
    }

    if user.Role != "admin" {
        return errors.New("permission denied")
    }

    return nil
}
```

### Summary

- Use `ng.Store` to attach values to the context.
- Use `ng.Load` to retrieve values from the context.
- Use `ng.MustLoad` when you are certain the value exists, which will panic if it does not.
- Use `ng.LoadOrStore` to retrieve a value or store a default if it does not exist.
- These utilities help maintain clean and modular code by avoiding global variables and passing data explicitly through the context.

---

## License

This project is licensed under the MIT License.
