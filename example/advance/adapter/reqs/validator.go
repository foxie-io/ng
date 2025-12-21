package reqs

import (
	"context"
	"fmt"
	"strings"

	"github.com/creasty/defaults"
	"github.com/foxie-io/ng"
	nghttp "github.com/foxie-io/ng/http"
	"github.com/go-playground/validator/v10"
)

var (
	valid = validator.New()
)

/*
Validate struct fields tagged with `validate` tag

For mode tags and usage, refer to https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Usage_and_Tags

Example:

	type CreateUserRequest struct {
		Name  string `json:"name" validate:"required,min=3,max=32"`
		Email string `json:"email" validate:"required,email"`
	}
*/
func Validate(s any) ng.Handler {
	return func(ctx context.Context) error {
		if errs := valid.Struct(s); errs != nil {
			verrs := errs.(validator.ValidationErrors)

			validateError := map[string]string{}
			for _, verr := range verrs {
				fieldName := verr.Field()
				prefix, after := fieldName[:1], fieldName[1:]
				fieldName = fmt.Sprintf("%s%s", strings.ToLower(prefix), after)
				validateError[fieldName] = verr.Tag()
			}

			return nghttp.NewErrInvalidArgument().Update(nghttp.Meta("error", validateError))
		}
		return nil
	}
}

/*
Set default values to struct fields tagged with `default` tag
Example:

	type CreateUserRequest struct {
		Name     string `json:"name" default:"john"`
		Age      int    `json:"age" default:"18"`
		Verified bool   `json:"verified" default:"false"`
	}
*/
func SetDefault(dest any) ng.Handler {
	return func(ctx context.Context) error {
		return defaults.Set(dest)
	}
}
