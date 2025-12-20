package ng

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func normolizePath(path string) string {
	if path == "/" {
		return ""
	}

	path = strings.TrimSuffix(path, "/")
	path = strings.TrimPrefix(path, "/")
	return "/" + path
}

func ExtractControllerRoutes(app App, config ControllerInitializer) ([]Route, error) {
	routes := []Route{}
	configType := reflect.TypeOf(config)
	configValue := reflect.ValueOf(config)

	for i := 0; i < configType.NumMethod(); i++ {
		funcName := configType.Method(i)

		if funcName.Type.NumOut() == 1 {
			meththodValue := configValue.MethodByName(funcName.Name)
			routeFn, ok := meththodValue.Interface().(func() Route)
			if !ok {
				continue
			}

			route, ok := routeFn().(*route)
			if !ok {
				return nil, fmt.Errorf("failed to build route")
			}

			route.name = strings.Replace(fmt.Sprintf("%T.%s", config, funcName.Name), "*", "", 1)
			route.core.preExecutes = append(route.core.preExecutes,
				func(ctx context.Context) {
					GetContext(ctx).setOwner(app, config, route)
				},
			)

			routes = append(routes, route)
		}
	}

	return routes, nil
}
