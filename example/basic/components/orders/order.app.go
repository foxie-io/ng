package orders

import (
	"example/basic/middlewares/limiter"
	"time"

	"github.com/foxie-io/ng"
)

func NewOrderApp() ng.App {
	orderApp := ng.NewApp(
		ng.WithPrefix("/order-app"),
		limiter.WithRateLimit(&limiter.Config{
			Limit:  10,
			Window: time.Minute,
		}),
	)

	orderApp.AddController(NewOrderController())
	return orderApp
}
