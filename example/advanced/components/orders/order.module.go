package orders

import (
	"example/advanced/router"

	"go.uber.org/fx"
)

var Module = fx.Module("orders",
	fx.Provide(
		NewOrderService,
	),
	router.GlobalController.Add(NewOrderController),
)
