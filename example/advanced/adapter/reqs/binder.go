package reqs

import (
	"context"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
)

var (
	binder = echo.DefaultBinder{}
)

/*
	Bind all (Path, Query, Body) to struct

Example:

	// URL: /users/:id?active=true
	// JSON Body: { "name": "john", "email": ""}
	// Swagger Decoration:
	// @Param _ path dtos.UpdateUserRequest true "Path Params"
	// @Param _ query dtos.UpdateUserRequest true "Query Params"
	// @Param _ body dtos.UpdateUserRequest true "Body"
	// @Router /users/{id} [put]
	type UpdateUserRequest struct {
		ID     string `param:"id"`
		Active bool   `query:"active"`
		Name   string `json:"name" form:"name" xml:"name"`
		Email  string `json:"email" form:"email" xml:"email"`
	}
*/
func BindAuto(dest interface{}) ng.Handler {
	return func(ctx context.Context) error {
		ectx := ng.MustLoad[echo.Context](ctx)
		if err := binder.Bind(dest, ectx); err != nil {
			return err
		}
		return nil
	}
}

/*
Bind Request Body to struct

Example:

	// JSON Body: { "name": "john", "email": ""}
	// Swagger Decoration:
	// @Param _ body dtos.CreateUserRequest true "Body"
	// @Router /users [post]
	type CreateUserRequest struct {
		Name  string `json:"name" form:"name" xml:"name"`
		Email string `json:"email" form:"email" xml:"email"`
	}
*/
func BindBody(dest interface{}) ng.Handler {
	return func(ctx context.Context) error {
		ectx := ng.MustLoad[echo.Context](ctx)
		if err := binder.BindBody(ectx, dest); err != nil {
			return nghttp.NewErrInvalidArgument().Update(
				nghttp.Meta("error", err.Error()),
			)
		}
		return nil
	}
}

/*
Bind Path Parameters to struct

Swagger: //

Example:

	// URL: /users/:id
	// Swagger Decoration:
	// @Param _ path dtos.UserPathParams true "Path Params"
	// @Router /users/{id} [get]
	type UserPathParams struct {
		ID string `param:"id"`
	}
*/
func BindParam(dest interface{}) ng.Handler {
	return func(ctx context.Context) error {
		ectx := ng.MustLoad[echo.Context](ctx)
		if err := binder.BindPathParams(ectx, dest); err != nil {
			return nghttp.NewErrInvalidArgument().Update(
				nghttp.Meta("error", err.Error()),
			)
		}
		return nil
	}
}

/*
Bind URL Query Parameters to struct

Example:

	// URL: ?name=john&email=join@gmail.com"
	// Swagger Decoration:
	// @Param _ query dtos.Filter true "Query Params"
	// @Router /users [get]
	type Filter struct {
		Name  string `query:"name"`
		Email string `query:"email"`
	}
*/
func BindQuery(dest interface{}) ng.Handler {
	return func(ctx context.Context) error {
		ectx := ng.MustLoad[echo.Context](ctx)
		if err := binder.BindQueryParams(ectx, dest); err != nil {
			return nghttp.NewErrInvalidArgument().Update(
				nghttp.Meta("error", err.Error()),
			)
		}
		return nil
	}
}
