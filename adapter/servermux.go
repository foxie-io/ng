package ngadapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/foxie-io/ng"
)

func ServeMuxResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
	w := ng.MustLoad[http.ResponseWriter](ctx)
	if info.HttpResponse != nil {
		w.WriteHeader(info.HttpResponse.StatusCode())
		bytes, _ := json.Marshal(info.HttpResponse.Response())
		_, _ = w.Write(bytes)
		return nil
	}

	log.Printf("no http response found in response info: raw:%v, stacks:%v", info.Raw, string(info.Stack))
	status := http.StatusInternalServerError
	http.Error(w, http.StatusText(status), status)
	return nil
}

func ServeMuxHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := ng.AcquireContext(r.Context())

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

func ServeMuxRegisterRoutes(ng ng.App, mux *http.ServeMux) {
	for _, route := range ng.Routes() {
		mux.HandleFunc(route.Path(), ServeMuxHandler(route.Handler))
	}
}
