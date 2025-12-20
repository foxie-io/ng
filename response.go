package ng

import (
	"context"

	nghttp "github.com/foxie-io/ng/http"
)

func Respond(ctx context.Context, val nghttp.HttpResponse) error {
	return GetContext(ctx).SetResponse(val).Response()
}

func ResponseAny(ctx context.Context, val any) error {
	switch t := val.(type) {
	case nghttp.HttpResponse:
		return GetContext(ctx).SetResponse(t).Response()
	case error:
		return GetContext(ctx).SetResponse(nghttp.WrapError(t)).Response()
	default:
		panic(val)
	}
}
