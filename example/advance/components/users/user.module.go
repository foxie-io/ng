package users

import (
	"example/advance/router"

	"go.uber.org/fx"
)

var Module = fx.Module("users",
	fx.Provide(
		NewUserService,
	),
	router.GlobalController.Add(NewUserController),
)
