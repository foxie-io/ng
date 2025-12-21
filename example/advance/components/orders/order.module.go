package orders

import (
	"example/advance/router"

	"go.uber.org/fx"
)

var Module = fx.Module("orders",
	fx.Provide(
		NewOrderService,
	),
	router.GlobalController.Add(NewOrderController),
)
