package ng

var _ Controller = (*controller)(nil)

type (
	// Controller is the main interface for handling routes
	Controller interface {
		Core() Core
		Routes() []Route
	}

	controller struct {
		core   *core
		routes []Route
	}
)

// Routes get controller routes
func (c *controller) Routes() []Route {
	return c.routes
}

// Core get controller core
func (c *controller) Core() Core {
	return c.core
}

// NewController create new controller instance
func NewController(opts ...Option) Controller {
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

// ControllerInitializer interface for initializing controller
type ControllerInitializer interface {
	InitializeController() Controller
}

// DefaultControllerInitializer default implementation
type DefaultControllerInitializer struct {
}

// InitializeController default implementation
func (c *DefaultControllerInitializer) InitializeController() Controller {
	return NewController()
}
