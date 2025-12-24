package ng

import "context"

// Interceptor executes after guards and before the final handler.
//
// Interceptors can wrap the handler execution and are commonly used for:
//   - Response transformation
//   - Error handling
//   - Metrics and timing for private operations due to its position after guards
//   - Transaction management
//
// Interceptors must call next to continue the request flow.
/*
Intercept(ctx context.Context, next Handler) {
	// middleware can write like this too to have more control that cover guard's opertions

	// before operation
	start := time.Now()

	// use defer to guarantee execution, event if an panic or error occurs
	defer func(){
		// after operation

		delay := time.Since(start)
		log.Printf("operation took %s", delay)
	}()

	// current operation
	next(ctx)
}
*/
type Interceptor interface {
	Intercept(ctx context.Context, next Handler)
}

// InterceptorFunc is an adapter to allow the use of ordinary functions as Interceptors.
type InterceptorFunc func(ctx context.Context, next Handler)

// Intercept calls f(ctx, next).
func (ifunc InterceptorFunc) Intercept(ctx context.Context, next Handler) {
	ifunc(ctx, next)
}
