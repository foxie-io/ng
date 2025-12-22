package router

import (
	"fmt"

	"github.com/foxie-io/ng"
	"go.uber.org/fx"
)

type ControllerGroup string

var (
	GlobalController ControllerGroup = "global_controller_initializers"
)

func (group ControllerGroup) Add(controllerInitializer any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			controllerInitializer,
			fx.As(new(ng.ControllerInitializer)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, string(group))),
		),
	)
}
