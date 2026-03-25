package mid

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/foundation/logger"
	"github.com/OpenJenie/goserve/foundation/web"
)

func TestErrorsReturnsForbiddenForAuthorizationFailures(t *testing.T) {
	t.Parallel()

	log := logger.New(io.Discard, logger.LevelInfo, "TEST", func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	})

	handler := Errors(log)(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return auth.NewAuthErrorWithStatus(http.StatusForbidden, "forbidden")
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/exampleauth", nil)
	rec := httptest.NewRecorder()
	ctx := web.SetValues(req.Context(), &web.Values{})

	if err := handler(ctx, rec, req); err != nil {
		t.Fatalf("unexpected middleware error: %v", err)
	}

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d", http.StatusForbidden, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "Forbidden") {
		t.Fatalf("expected Forbidden response body, got %s", rec.Body.String())
	}
}
