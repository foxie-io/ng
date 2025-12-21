package users

import (
	"context"
	"example/advance/adapter/reqs"
	"example/advance/components/users/dtos"

	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type UserController struct {
	ng.DefaultControllerInitializer
	user_s *UserService
}

func NewUserController(user_s *UserService) *UserController {
	return &UserController{
		user_s: user_s,
	}
}

func (con *UserController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/users"),
	)
}

func (con *UserController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				body dtos.CreateUserRequest
			)

			return ng.Handle(
				reqs.BindBody(&body),
				reqs.Validate(&body),
				func(ctx context.Context) error {
					resp, err := con.user_s.CreateUser(body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *UserController) Get() ng.Route {
	return ng.NewRoute(http.MethodGet, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.GetUserRequest
			)
			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.user_s.GetUser(param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *UserController) GetAll() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				query dtos.ListUsersRequest
			)
			return ng.Handle(
				reqs.BindQuery(&query),
				reqs.SetDefault(&query),
				reqs.Validate(&query),
				func(ctx context.Context) error {
					resp := con.user_s.GetAllUsers(&query)
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *UserController) Update() ng.Route {
	return ng.NewRoute(http.MethodPut, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.DeleteUserRequest
				body  dtos.UpdateUserRequest
			)

			return ng.Handle(
				reqs.BindBody(&body), reqs.BindParam(&param),
				reqs.Validate(&body), reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.user_s.UpdateUser(param.ID, &body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

func (con *UserController) Delete() ng.Route {
	return ng.NewRoute(http.MethodDelete, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.DeleteUserRequest
			)
			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),
				func(ctx context.Context) error {
					resp := con.user_s.DeleteUser(&param)
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}
