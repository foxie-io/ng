package ng

var _ Controller = (*controller)(nil)

type (
	Controller interface {
		Core() Core
		Routes() []Route
	}

	controller struct {
		core   *core
		routes []Route
	}
)

func (c *controller) Routes() []Route {
	return c.routes
}

func (c *controller) Core() Core {
	return c.core
}

func NewController(opts ...option) Controller {
	controller := &controller{core: newCore()}

	config := newConfig()
	config.bindController(controller)
	config.bindCore(controller.core)
	config.update(opts...)

	return controller
}

func (c *controller) addRoute(routes ...Route) Controller {
	for _, r := range routes {
		r.(*route).addPreCore(c.core)
		c.routes = append(c.routes, r)
	}

	return c
}

func (c *controller) build(app *app, config ControllerInitializer) Controller {
	if c.core.built.Load() {
		panic("controller already built")
	}

	defer c.core.built.Store(true)
	routes, err := ExtractControllerRoutes(app, config)
	if err != nil {
		panic(err)
	}

	c.addRoute(routes...)
	return c
}

type ControllerInitializer interface {
	InitializeController() Controller
}

// ControllerInitializer
type DefaultControllerInitializer struct {
}

func (c *DefaultControllerInitializer) InitializeController() Controller {
	return NewController()
}
