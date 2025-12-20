package orders

import (
	"context"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type OrderTodoController struct {
	ng.DefaultControllerInitializer
}

func NewOrderController() *OrderTodoController {
	return &OrderTodoController{}
}

func (con *OrderTodoController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/order-todo-controller"),
	)
}

// curl 'localhost:8080/order-app/order-todo-controller/hello-world'
func (con *OrderTodoController) Hello() ng.Route {
	return ng.NewRoute(http.MethodGet, "/hello-world",
		ng.WithHandler(func(ctx context.Context) error {
			rc := ng.GetContext(ctx)
			txt := rc.Route().Name() + " Hello World!"
			return ng.Respond(ctx, nghttp.NewResponse(txt))
		}),
	)
}
