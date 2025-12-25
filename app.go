package ng

type (

	// App is the main application interface
	App interface {
		Core() Core
		Routes() []Route
		Build() App
		AddSubApp(app ...App)
		AddController(configs ...ControllerInitializer)
		AddRoute(routes ...Route)
	}

	app struct {
		core *core

		configs []ControllerInitializer

		routes []Route

		subApps []App
	}
)

func (a *app) Core() Core {
	return a.core
}

// NewApp creates a new App instance
func NewApp(opts ...Option) App {
	app := &app{core: newCore()}
	return app.update(opts...)
}

func (a *app) update(opts ...Option) App {
	if a.core.built.Load() {
		panic("app is builded can't update")
	}

	config := newConfig()
	config.bindApp(a)
	config.bindCore(a.core)
	config.update(opts...)
	return a
}

func (a *app) AddSubApp(subs ...App) {
	a.subApps = append(a.subApps, subs...)

}

func (a *app) AddController(configs ...ControllerInitializer) {
	if a.core.built.Load() {
		panic("app already built")
	}

	a.configs = append(a.configs, configs...)
}

func (a *app) Routes() []Route {
	return a.routes
}

func (a *app) AddRoute(routes ...Route) {
	for _, r := range routes {
		r.(*route).addPreCore(a.core)
		a.routes = append(a.routes, r)
	}
}

func (a *app) buildController() App {
	// extract routes from configs
	for _, config := range a.configs {
		controller := config.InitializeController().(*controller)
		controller.build(a, config)

		for _, r := range controller.Routes() {
			a.AddRoute(r)
		}
	}

	return a
}

func (a *app) buildRouter() App {
	for _, r := range a.routes {
		r.(*route).build()
	}

	return a
}

func (a *app) extractRouterFromSubApp() {
	for _, sub := range a.subApps {
		subApp := sub.(*app).buildController()
		for _, r := range subApp.Routes() {
			a.AddRoute(r)
		}
	}
}

func (a *app) Build() App {
	if a.core.built.Load() {
		panic("app already built")
	}

	defer a.core.built.Store(true)

	// extract routes from configs

	a.buildController()

	// add sub apps
	a.extractRouterFromSubApp()

	// build routes
	a.buildRouter()

	return a
}
