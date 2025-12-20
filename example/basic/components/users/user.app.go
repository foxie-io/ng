package users

import "github.com/foxie-io/ng"

func NewUserApp() ng.App {
	userApp := ng.NewApp(
		ng.WithPrefix("/user-app"),
	)

	userApp.AddController(NewUserController())
	return userApp
}
