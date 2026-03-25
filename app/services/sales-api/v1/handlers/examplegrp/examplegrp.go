package examplegrp

import (
	"context"
	"net/http"

	"github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/foundation/web"
)

// Handlers manages the example endpoints for the starter project.
type Handlers struct {
	build string
}

// New constructs a Handlers value for the example group.
func New(build string) *Handlers {
	return &Handlers{build: build}
}

// Example returns starter metadata without authentication requirements.
func (h *Handlers) Example(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Build  string `json:"build"`
	}{
		Name:   "goserve-starter",
		Status: "available",
		Build:  h.build,
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}

// ExampleAuth returns the authenticated claims for API integration smoke tests.
func (h *Handlers) ExampleAuth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims := auth.GetClaims(ctx)

	data := struct {
		Status  string   `json:"status"`
		Subject string   `json:"subject"`
		Roles   []string `json:"roles"`
		Build   string   `json:"build"`
	}{
		Status:  "authorized",
		Subject: claims.Subject,
		Roles:   claims.Roles,
		Build:   h.build,
	}

	return web.Respond(ctx, w, data, http.StatusOK)
}
