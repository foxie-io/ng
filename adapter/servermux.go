package ngadapter

// HTTP adapter for net/http ServeMux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
)

// ServeMuxResponseHandler write HTTPResponse to http.ResponseWriter
func ServeMuxResponseHandler(ctx context.Context, info nghttp.HTTPResponse) error {
	var (
		w     = ng.MustLoad[http.ResponseWriter](ctx)
		value []byte
	)

	switch v := info.(type) {
	case *nghttp.Response:
		w.Header().Set("content-type", "application/json")
		value, _ = json.Marshal(v.Response())

	case *nghttp.PanicError:
		fmt.Println("recieve (*nghttp.PanicError)", v.Value())
		value, _ = json.Marshal(info.Response())

	case *nghttp.RawResponse:
		value = v.Value()

	default:
		fmt.Println("unknown in response", info.Response())
		value = []byte("unknown in response value")
	}

	w.WriteHeader(info.StatusCode())
	_, _ = w.Write(value)
	return nil
}

// ServeMuxHandler create http.HandlerFunc from ng.Handler
func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := ng.NewContext(r.Context())

		// store in context
		ng.Store(ctx, w)
		ng.Store(ctx, r)

		// can extract from ctx if needed
		// w := ng.MustLoad[http.ResponseWriter](ctx)
		// r := ng.MustLoad[*http.Request](ctx)

		// invoke the handler
		scopeHandler()(ctx)
	}
}

// ServeMuxRegisterRoutes register all routes from ng.App into http.ServeMux
func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		muxPath := fmt.Sprintf("%s %s", route.Method(), route.Path())
		mux.HandleFunc(muxPath, ServeMuxHandler(route.Handler))
	}
}
