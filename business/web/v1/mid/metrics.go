package mid

import (
	"context"
	"net/http"

	metrics "github.com/OpenJenie/goserve/business/web/v1/mertrics"
	"github.com/OpenJenie/goserve/foundation/web"
)

// Metrics updates program counters.
func Metrics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)

			err := handler(ctx, w, r)

			n := metrics.AddRequests(ctx)
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}

		return h
	}

	return m
}
