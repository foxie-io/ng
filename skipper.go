package ng

import (
	"context"
	"fmt"
)

type (
	// ID identifies a component (middleware, guard, interceptor)
	// that can be conditionally skipped during request execution.
	//
	// The ID must be stable and unique for the component type.
	ID interface {
		NgID() string
	}

	// skipperKey is an internal metadata key used to store skipper IDs
	// at the route or handler level.
	skipperKey struct{}

	// DefaultID allows skipping middleware, guards, or interceptors
	// by their concrete type.
	//
	// Example:
	//
	//	ng.WithSkipper(ng.DefaultID[AuthGuard]{})
	//
	// This will prevent AuthGuard from executing for the configured route.
	DefaultID[T any] struct{}
)

// NgID returns a unique identifier for the given generic type T.
func (s DefaultID[T]) NgID() string {
	return fmt.Sprintf("skipper_%T", s)
}

// WithSkip attaches skipper metadata to a route or handler.
//
// Any middleware, guard, or interceptor implementing Skipper
// whose SkipID matches one of the provided skippers will be skipped
// during request execution.
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
