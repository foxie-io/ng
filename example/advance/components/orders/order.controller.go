package orders

import (
	"context"
	"example/advance/adapter/reqs"
	"example/advance/components/orders/dtos"

	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

type OrderController struct {
	ng.DefaultControllerInitializer
	order_s *OrderService
}

func NewOrderController(order_s *OrderService) *OrderController {
	return &OrderController{
		order_s: order_s,
	}
}

func (con *OrderController) InitializeController() ng.Controller {
	return ng.NewController(
		ng.WithPrefix("/orders"),
	)
}

// @Tags        Orders
// @Summary     Order create
// @ID          order-create
// @Accept      json
// @Produce     json
// @Param       state body dtos.CreateOrderRequest true "Create Order Request"
// @Success     200 {object} nghttp.Response{data=dtos.CreateOrderResponse}
// @Router      /orders [post]
func (con *OrderController) Create() ng.Route {
	return ng.NewRoute(http.MethodPost, "/",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				body dtos.CreateOrderRequest
			)

			return ng.Handle(
				reqs.BindBody(&body),
				reqs.Validate(&body),

				func(ctx context.Context) error {
					resp, err := con.order_s.CreateOrder(ctx, body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        Orders
// @Summary     GetOrder orders
// @ID          get-orders
// @Accept      json
// @Produce     json
// @Param       struct path dtos.PathID true  "Order ID"
// @Success     200 {object} nghttp.Response{data=dtos.GetOrderResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router /orders/{id} [get]
func (con *OrderController) GetOrder() ng.Route {
	return ng.NewRoute(http.MethodGet, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
			)

			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),

				func(ctx context.Context) error {
					resp, err := con.order_s.GetOrder(ctx, param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        Orders
// @Summary     Get all orders
// @ID          get-all-orders
// @Accept      json
// @Produce     json
// @Param       query query dtos.ListOrdersRequest true "List Orders Request"
// @Success     200 {object} nghttp.Response{data=[]dtos.GetAllOrdersResponse}
// @Router      /orders [get]
func (con *OrderController) GetOrders() ng.Route {
	return ng.NewRoute(http.MethodGet, "/",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				query dtos.ListOrdersRequest
			)
			return ng.Handle(
				reqs.BindQuery(&query),
				reqs.SetDefault(&query),
				reqs.Validate(&query),

				func(ctx context.Context) error {
					resp := con.order_s.GetOrders(ctx, &query)
					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags        Orders
// @Summary     Update an order
// @ID          update-order
// @Accept      json
// @Produce     json
// @Param       id path dtos.PathID true "Order ID"
// @Param       body body dtos.UpdateOrderRequest true "Update Order Request"
// @Success     200 {object} nghttp.Response{data=dtos.UpdateOrderResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router      /orders/{id} [put]
func (con *OrderController) Update() ng.Route {
	return ng.NewRoute(http.MethodPut, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
				body  dtos.UpdateOrderRequest
			)

			return ng.Handle(
				reqs.BindParam(&param), reqs.BindBody(&body),
				reqs.Validate(&param), reqs.Validate(&body),
				func(ctx context.Context) error {
					resp, err := con.order_s.UpdateOrder(ctx, param.ID, &body)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}

// @Tags 		Orders
// @Summary 	Delete an order
// @Accept 		json
// @Produce 	json
// @Param 		state path dtos.PathID true "Order ID"
// @Success 	200 {object} nghttp.Response{data=dtos.DeleteOrderResponse}
// @Failure     404 {object} nghttp.Response "NOT_FOUND"
// @Router /orders/{id} [delete]
func (con *OrderController) Delete() ng.Route {
	return ng.NewRoute(http.MethodDelete, "/:id",
		ng.WithScopeHandler(func() ng.Handler {
			var (
				param dtos.PathID
			)
			return ng.Handle(
				reqs.BindParam(&param),
				reqs.Validate(&param),
				func(ctx context.Context) error {
					resp, err := con.order_s.DeleteOrder(ctx, param.ID)
					if err != nil {
						return err
					}

					return ng.Respond(ctx, nghttp.NewResponse(resp))
				},
			)
		}),
	)
}
