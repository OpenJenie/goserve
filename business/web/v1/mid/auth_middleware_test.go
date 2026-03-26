package mid_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/business/web/v1/mid"
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
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusOK {
				t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
			}
		})

		t.Run("unauthorized-missing-header", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusUnauthorized {
				t.Errorf("got status %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})

		t.Run("unauthorized-empty-bearer", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer ")
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusUnauthorized {
				t.Errorf("got status %d, want %d", w.Code, http.StatusUnauthorized)
			}
		})

		t.Run("unauthorized-bad-token", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.Header.Set("Authorization", "Bearer bad-token")
			w := httptest.NewRecorder()

			if err := handler(context.Background(), w, r); err != nil {
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusUnauthorized {
				t.Errorf("got status %d, want %d", w.Code, http.StatusUnauthorized)
			}
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
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusOK {
				t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
			}
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
				t.Fatalf("handler: %s", err)
			}

			if w.Code != http.StatusForbidden {
				t.Errorf("got status %d, want %d", w.Code, http.StatusForbidden)
			}
		})
	}

	test(t)
}

// =============================================================================

func testKeyStore(t *testing.T) (auth.KeyLookup, string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating key: %v", err)
	}

	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privatePEM := string(pem.EncodeToMemory(privateBlock))

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshaling public key: %v", err)
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	publicPEM := string(pem.EncodeToMemory(publicBlock))

	kid := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	ks := &simpleKeyStore{
		store: map[string]string{
			kid: privatePEM,
		},
		pubStore: map[string]string{
			kid: publicPEM,
		},
	}
	return ks, kid
}

type simpleKeyStore struct {
	store    map[string]string
	pubStore map[string]string
}

func (ks *simpleKeyStore) PrivateKey(kid string) (string, error) {
	if pem, ok := ks.store[kid]; ok {
		return pem, nil
	}
	return "", errors.New("not found")
}

func (ks *simpleKeyStore) PublicKey(kid string) (string, error) {
	if pem, ok := ks.pubStore[kid]; ok {
		return pem, nil
	}
	return "", errors.New("not found")
}
