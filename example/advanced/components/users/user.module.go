package users

import (
	"example/advanced/router"

	"go.uber.org/fx"
)

var Module = fx.Module("users",
	fx.Provide(
		NewUserService,
	),
	router.GlobalController.Add(NewUserController),
)
