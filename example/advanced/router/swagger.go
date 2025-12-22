package router

import (
	"context"
	"fmt"
	"time"

	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/swag"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

var (
	// implment checker
	_ interface {
		ng.ControllerInitializer
		nghttp.HttpResponse
	} = (*SwaggerDocs)(nil)
)

type SwaggerDocs struct {
	ng.DefaultControllerInitializer
	html     string
	dir      string
	filename string
}

func NewSwaggerUI(spec *swag.Spec) *SwaggerDocs {
	context := swag.GetSwagger(spec.InstanceName()).ReadDoc()

	su := &SwaggerDocs{
		dir:      "./docs",
		filename: "swagger.json",
	}

	// scalar: swagger client ui
	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecContent: context,
		IsEditable:  true,
	})

	if err != nil {
		panic(fmt.Errorf("failed to generate swagger ui: %w", err))
	}

	su.html = html
	return su
}

// nghttp.HttpResponse implementation for ng.Respond(ctx, su)
func (su *SwaggerDocs) StatusCode() int { return 200 }
func (su *SwaggerDocs) Response() any   { return su.html }

func (su *SwaggerDocs) UI() ng.Route {

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Swagger UI available at:", "http://localhost:8080/docs")
	}()

	return ng.NewRoute("GET", "/docs",
		ng.WithHandler(func(ctx context.Context) error {
			ng.Respond(ctx, su)
			return nil
		}),
		// custom response handler to serve HTML
		ng.WithResponseHandler(func(ctx context.Context, info nghttp.HttpResponse) error {
			c := ng.MustLoad[echo.Context](ctx)

			resp, ok := info.(*SwaggerDocs)
			if !ok {
				return c.String(500, "invalid response type")
			}

			c.HTML(resp.StatusCode(), resp.Response().(string))
			return nil
		}),
	)
}
