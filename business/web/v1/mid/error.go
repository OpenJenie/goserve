package mid

import (
	"context"
	"net/http"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/business/web/v1/response"
	"github.com/OpenJenie/goserve/foundation/logger"
	"github.com/OpenJenie/goserve/foundation/validate"
	"github.com/OpenJenie/goserve/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				var er response.ErrorDocument
				var status int

				switch {
				case response.IsError(err):
					reqErr := response.GetError(err)

					if validate.IsFieldErrors(reqErr.Err) {
						fieldErrors := validate.GetFieldErrors(reqErr.Err)
						er = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = reqErr.Status
						break
					}

					er = response.ErrorDocument{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				case auth.IsAuthError(err):
					authErr := auth.GetAuthError(err)
					status = http.StatusUnauthorized
					if authErr != nil && authErr.Status() != 0 {
						status = authErr.Status()
					}
					er = response.ErrorDocument{
						Error: http.StatusText(status),
					}

				default:
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
