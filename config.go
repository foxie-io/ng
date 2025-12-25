package ng

// allow option to update universal app,core,route,controller etc...
type config struct {
	app        *app
	core       *core
	route      *route
	controller *controller
}

// mutate config.route affect route
func (c *config) bindRoute(route *route) {
	c.route = route
}

// mutate config.core affect core
func (c *config) bindCore(core *core) {
	c.core = core
}

// mutate config.controller will affect controller
func (c *config) bindController(controller *controller) {
	c.controller = controller
}

// mutate config.app will affect app
func (c *config) bindApp(app *app) {
	c.app = app
}

// to bind before update
func (c *config) update(opts ...Option) {
	for _, u := range opts {
		u(c)
	}
}

// Option is ng configuration option
type Option func(*config)

// Opitons combines multiple options into one
func Opitons(opts ...Option) Option {
	return func(c *config) {
		for _, o := range opts {
			o(c)
		}
	}
}

func newConfig() *config {
	return &config{}
}
