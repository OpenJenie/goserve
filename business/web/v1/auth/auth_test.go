package auth

import (
	"context"
	"testing"
	"time"

	"github.com/OpenJenie/goserve/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	ks, kid := testKeyStore(t)
	authn, err := New(Config{
		KeyLookup: ks,
		Issuer:    "test-issuer",
	})
	if err != nil {
		t.Fatalf("new auth: %v", err)
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    "test-issuer",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
		Roles: []string{"ADMIN"},
	}

	token, err := authn.GenerateToken(kid, claims)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	got, err := authn.Authenticate(context.Background(), "Bearer "+token)
	if err != nil {
		t.Fatalf("authenticate: %v", err)
	}

	if got.Subject != claims.Subject {
		t.Fatalf("expected subject %q, got %q", claims.Subject, got.Subject)
	}
}

func TestAuthenticateRejectsWrongIssuer(t *testing.T) {
	t.Parallel()

	ks, kid := testKeyStore(t)
	authn, err := New(Config{
		KeyLookup: ks,
		Issuer:    "expected-issuer",
	})
	if err != nil {
		t.Fatalf("new auth: %v", err)
	}

	token, err := authn.GenerateToken(kid, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    "wrong-issuer",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
		Roles: []string{"ADMIN"},
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	if _, err := authn.Authenticate(context.Background(), "Bearer "+token); err == nil {
		t.Fatal("expected issuer mismatch to fail authentication")
	}
}

func testKeyStore(t *testing.T) (*keystore.KeyStore, string) {
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
