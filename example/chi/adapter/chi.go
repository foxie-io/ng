package adapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/foxie-io/ng"
	"github.com/go-chi/chi/v5"
)

func ChiResponseHandler(ctx context.Context, info *ng.ResponseInfo) error {
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

func ChiHandler(scopeHandler func() ng.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := ng.AcquireContext(r.Context())

		// store http request and response writer
		ng.Store(ctx, w)
		ng.Store(ctx, r)

		// get http request and response writer from ng ctx
		// w := ng.MustLoad[http.ResponseWriter](ctx)
		// r := ng.MustLoad[*http.Request](ctx)
		_ = scopeHandler()(ctx)
	}
}

func ChiRegisterRoutes(ng ng.App, router chi.Router) {
	for _, route := range ng.Routes() {
		chiHandler := ChiHandler(route.Handler)
		router.Method(route.Method(), route.Path(), chiHandler)
	}
}
