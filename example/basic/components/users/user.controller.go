package users

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type UserTodoController struct {
	ng.DefaultControllerInitializer
}

func NewUserController() *UserTodoController {
	return &UserTodoController{}
}

func (con *UserTodoController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/user-todo-controller"),
	)
}

func (con *UserTodoController) Hello() ng.Route {
	return ng.NewRoute(http.MethodGet, "/hello-world",
		ng.WithHandler(func(ctx context.Context) error {
			rc := ng.GetContext(ctx)
			txt := rc.Route().Name() + " Hello World!"
			return ng.Respond(ctx, nghttp.NewResponse(txt))
		}),
	)
}
