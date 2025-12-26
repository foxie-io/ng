package ng

// Skipper provides a way to conditionally skip middleware, guards, or interceptors
// during request execution.
//
// By implementing the Skipper interface and returning a unique NgID,
// features can be identified for skipping based on route or handler metadata.
//
// This is useful for scenarios such as:
//   - Bypassing authentication guards for public endpoints
//   - Skipping logging middleware for health check routes
//   - Omitting response transformation interceptors for specific handlers
/// Example:
/*
type AuthGuard struct {
	// ...
}

func (ag *AuthGuard) NgID() string {
	return "auth_guard"
}

func (ag *AuthGuard) Allow(ctx context.Context) error {
	// guard logic
}
*/

import (
	"context"
	"fmt"
)

type (
	// ID represents a unique identifier for skippable features.
	ID interface {
		NgID() string
	}

	// skipperKey is used as a metadata key for storing skipper IDs.
	skipperKey struct{}

	// DefaultID is a generic implementation of the ID interface.
	// example:
	/*
		// allow guard to be skipped by using DefaultID[type]
		type AuthGuard struct{
			DefaultID[AuthSkipper]
		}

		// allow middleware to be skipped by using DefaultID[type]
		type LoggingMiddleware struct{
			DefaultID[LoggingSkipper]
		}

		// allow interceptor to be skipped by using DefaultID[type]
		type ResponseInterceptor struct{
			DefaultID[ResponseSkipper]
		}
	*/
	DefaultID[T any] struct{}
)

// NgID returns a unique identifier for the given generic type T.
func (s DefaultID[T]) NgID() string {
	return fmt.Sprintf("skipper_%T", s)
}

// WithSkip can be used in app, controller, route to skip certain skippable features.
// For example, to skip certain guards, middlewares, or interceptors.
func WithSkip(skippers ...ID) Option {
	skipIds := make([]string, len(skippers))
	for i := range skippers {
		skipIds[i] = skippers[i].NgID()
	}

	return WithMetadata(skipperKey{}, skipIds)
}

const allGuard string = "all_guards"

// SkipAllGuards skips execution of all guards for the route or handler.
//
// This is useful for public endpoints such as health checks
// or authentication callbacks.
func SkipAllGuards() Option {
	return WithMetadata(skipperKey{}, []string{allGuard})
}

// getSkipperIds retrieves skipper IDs from the current request context.
// It returns nil if no skipper metadata is defined.
func getSkipperIds(ctx context.Context) []string {
	val, exists := GetContext(ctx).Route().Core().Metadata(skipperKey{})
	if !exists {
		return nil
	}

	return val.([]string)
}

// canSkip reports whether the given value should be skipped
// based on the provided skipper IDs.
func canSkip(val any, skipIds []string) bool {
	skipper, ok := val.(ID)
	if !ok {
		return false
	}

	for _, id := range skipIds {
		if id == skipper.NgID() {
			return true
		}
	}

	return false
}
