package users

import (
	"context"
	"example/advanced/adapter/reqs"
	"example/advanced/components/users/dtos"

	"net/http"

	_ "github.com/foxie-io/gormqs"
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

// @Tags        users
// @Summary     create an user
// @ID          user-create
// @Accept      json
// @Produce     json
// @Param       state body dtos.CreateUserRequest true "create user request"
// @Success     200 {object} nghttp.Response{data=dtos.CreateUserResponse}
// @Router      /users [post]
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
					resp, err := con.user_s.CreateUser(ctx, body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        users
// @Summary     get an user
// @ID          get-user
// @Accept      json
// @Produce     json
// @Param       userId path dtos.PathID true  "userId"
// @Success     200 {object} nghttp.Response{data=dtos.GetUserResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router /users/{id} [get]
func (con *UserController) Get() ng.Route {
	return ng.NewRoute(http.MethodGet, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
			)
			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.user_s.GetUser(ctx, param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        users
// @Summary     get users
// @ID          get-users
// @Accept      json
// @Produce     json
// @Param       query query dtos.ListUsersRequest true "List Orders Request"
// @Success     200 {object} nghttp.Response{data=gormqs.ListResulter[dtos.GetUserResponse]}
// @Router      /users [get]
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
					resp, err := con.user_s.GetAllUsers(ctx, &query)
					if err != nil {
						return err
					}
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        users
// @Summary     update an user
// @ID          update-user
// @Accept      json
// @Produce     json
// @Param       userId path dtos.PathID true "user ID"
// @Param       body body dtos.UpdateUserRequest true "Update user Request"
// @Success     200 {object} nghttp.Response{data=dtos.UpdateUserResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router      /users/{id} [put]
func (con *UserController) Update() ng.Route {
	return ng.NewRoute(http.MethodPut, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
				body  dtos.UpdateUserRequest
			)
			return ng.Handle(
				reqs.BindBody(&body), reqs.BindParam(&param),
				reqs.Validate(&body), reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.user_s.UpdateUser(ctx, param.ID, &body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags 		users
// @Summary 	delete an user
// @Accept 		json
// @Produce 	json
// @Param 		userId path dtos.PathID true "user ID"
// @Success 	200 {object} nghttp.Response{data=dtos.DeleteUserResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router /users/{id} [delete]
func (con *UserController) Delete() ng.Route {
	return ng.NewRoute(http.MethodDelete, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
			)
			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.user_s.DeleteUser(ctx, param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}
