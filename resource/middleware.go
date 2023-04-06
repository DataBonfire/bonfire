package resource

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

var StorageMiddleware = func(next middleware.Handler) middleware.Handler {
	return (middleware.Handler)(func(ctx context.Context, req interface{}) (interface{}, error) {
		ctx = context.WithValue(ctx, "storage", storage)
		return next(ctx, req)
	})
}

func MakeCors() func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"}),
		handlers.AllowedHeaders([]string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}),
		handlers.ExposedHeaders([]string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Cache-Control", "Content-Language", "Content-Type"}),
		handlers.AllowCredentials(),
	)
}

type pbvalidator interface {
	ValidateAll() error
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(pbvalidator); ok {
				if err := v.ValidateAll(); err != nil {
					// invalid RegisterRequest.Email: value must be a valid email address | caused by: mail: missing '@' or angle-addr; invalid RegisterRequest.Password: value length must be at least 8 runes
					return nil, errors.BadRequest("VALIDATOR", "").WithCause(err).WithMetadata(ToMetadata(err.Error()))
				}
			}
			return handler(ctx, req)
		}
	}
}

var errFinder = regexp.MustCompile(`invalid\s+[^.]+\.([^:]+):\s+([^;|]+)`)

func ToMetadata(msg string) map[string]string {
	finds := errFinder.FindAllStringSubmatch(msg, -1)
	metadata := map[string]string{}
	for _, find := range finds {
		//fmt.Println(ToSnakeCase(find[1]), find[2])
		//json.NewEncoder(os.Stdout).Encode(find)
		find[2] = strings.TrimSpace(find[2])
		if find[2] == "embedded message failed validation" {
			continue
		}
		metadata[ToSnakeCase(find[1])] = find[2]
	}
	return metadata
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
