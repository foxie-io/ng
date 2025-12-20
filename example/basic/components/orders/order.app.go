package orders

import "github.com/foxie-io/ng"

func NewOrderApp() ng.App {
	orderApp := ng.NewApp(
		ng.WithPrefix("/order-app"),
	)

	orderApp.AddController(NewOrderController())
	return orderApp
}
