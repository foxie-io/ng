package ng

import "context"

// Guard is used to determine if a request is allowed to proceed.
// Guards are typically used for authorization and access control.
//
// If a guard denies access, it should return an error indicating the reason for denial.
// If access is allowed, it should return nil.
/*
type RoleGuard struct {
	RequiredRole string
}

func (rg *RoleGuard) Allow(ctx context.Context) error {
	user := ng.MustLoad[User](ctx)
	if !user.HasRole(rg.RequiredRole) {
		return nghttp.NewErrPermissionDenied()
	}
	return nil
}
*/
type Guard interface {
	Allow(ctx context.Context) error
}

// GuardFunc is an adapter to allow the use of ordinary functions as Guards.
type GuardFunc func(ctx context.Context) error

// Allow calls f(ctx).
func (gf GuardFunc) Allow(ctx context.Context) error {
	return gf(ctx)
}
