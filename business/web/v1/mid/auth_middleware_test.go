package mid_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/business/web/v1/mid"
	"github.com/OpenJenie/goserve/foundation/keystore"
	"github.com/OpenJenie/goserve/foundation/logger"
	"github.com/OpenJenie/goserve/foundation/web"
	"github.com/golang-jwt/jwt/v4"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	test := func(t *testing.T) {
		ks, kid := testKeyStore(t)
		log := logger.New(io.Discard, logger.LevelInfo, "TEST", func(ctx context.Context) string {
			return web.GetTraceID(ctx)
		})

		a, err := auth.New(auth.Config{
			Log:       log,
			KeyLookup: ks,
			Issuer:    "test-issuer",
		})
		if err != nil {
			t.Fatalf("auth new: %s", err)
		}

		handler := mid.Errors(log)(mid.Authenticate(a)(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return web.Respond(ctx, w, "OK", http.StatusOK)
		}))

		t.Run("authorized", func(t *testing.T) {
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "user-123",
					Issuer:    "test-issuer",
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
				Roles: []string{"ADMIN"},
			}

			token, err := a.GenerateToken(kid, claims)
			if err != nil {
				t.Fatalf("generate token: %s", err)
			}

			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusOK, `"OK"`)
		})

		t.Run("unauthorized-missing-header", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
		})

		t.Run("unauthorized-empty-bearer", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer ")
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
		})

		t.Run("unauthorized-bad-token", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer bad-token")
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusUnauthorized, `{"error":"Unauthorized"}`)
		})
	}

	test(t)
}

func TestAuthorize(t *testing.T) {
	t.Parallel()

	test := func(t *testing.T) {
		ks, kid := testKeyStore(t)
		log := logger.New(io.Discard, logger.LevelInfo, "TEST", func(ctx context.Context) string {
			return web.GetTraceID(ctx)
		})

		a, err := auth.New(auth.Config{
			Log:       log,
			KeyLookup: ks,
			Issuer:    "test-issuer",
		})
		if err != nil {
			t.Fatalf("auth new: %s", err)
		}

		handler := mid.Errors(log)(mid.Authenticate(a)(mid.Authorize(a, auth.RuleAdminOnly)(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return web.Respond(ctx, w, "OK", http.StatusOK)
		})))

		t.Run("authorized", func(t *testing.T) {
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "admin-user",
					Issuer:    "test-issuer",
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
				Roles: []string{"ADMIN"},
			}

			token, err := a.GenerateToken(kid, claims)
			if err != nil {
				t.Fatalf("generate token: %s", err)
			}

			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusOK, `"OK"`)
		})

		t.Run("forbidden", func(t *testing.T) {
			claims := auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "user-only",
					Issuer:    "test-issuer",
					IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
				},
				Roles: []string{"USER"},
			}

			token, err := a.GenerateToken(kid, claims)
			if err != nil {
				t.Fatalf("generate token: %s", err)
			}

			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("unexpected handler error: %s", err)
			}

			assertResponse(t, w, http.StatusForbidden, `{"error":"Forbidden"}`)
		})
	}

	test(t)
}

// =============================================================================

func testKeyStore(t *testing.T) (auth.KeyLookup, string) {
	t.Helper()

	privatePEM, _, err := keystore.GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("generate key pair: %v", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		t.Fatalf("parse private key: %v", err)
	}

	const kid = "test-kid"
	return keystore.NewMap(map[string]keystore.PrivateKey{
		kid: {
			PK:  privateKey,
			PEM: privatePEM,
		},
	}), kid
}

func assertResponse(t *testing.T, w *httptest.ResponseRecorder, wantStatus int, wantBody string) {
	t.Helper()

	if w.Code != wantStatus {
		t.Fatalf("got status %d, want %d", w.Code, wantStatus)
	}

	if got := w.Body.String(); got != wantBody {
		t.Fatalf("got body %q, want %q", got, wantBody)
	}
}
