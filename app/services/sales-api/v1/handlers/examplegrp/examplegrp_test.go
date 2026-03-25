package examplegrp

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/foundation/web"
	"github.com/golang-jwt/jwt/v4"
)

func TestExample(t *testing.T) {
	t.Parallel()

	hdl := New("test-build")
	req := httptest.NewRequest(http.MethodGet, "/v1/example", nil)
	rec := httptest.NewRecorder()
	ctx := web.SetValues(req.Context(), &web.Values{})

	if err := hdl.Example(ctx, rec, req); err != nil {
		t.Fatalf("example handler: %v", err)
	}

	body := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(body, `"name":"goserve-starter"`) {
		t.Fatalf("expected starter name in response, got %s", body)
	}
}

func TestExampleAuth(t *testing.T) {
	t.Parallel()

	hdl := New("test-build")
	req := httptest.NewRequest(http.MethodGet, "/v1/exampleauth", nil)
	rec := httptest.NewRecorder()
	ctx := web.SetValues(req.Context(), &web.Values{})
	ctx = auth.SetClaims(ctx, auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "user-123",
		},
		Roles: []string{"ADMIN"},
	})

	if err := hdl.ExampleAuth(ctx, rec, req); err != nil {
		t.Fatalf("example auth handler: %v", err)
	}

	body := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
	if !strings.Contains(body, `"subject":"user-123"`) {
		t.Fatalf("expected subject in response, got %s", body)
	}
	if !strings.Contains(body, `"roles":["ADMIN"]`) {
		t.Fatalf("expected roles in response, got %s", body)
	}
}
