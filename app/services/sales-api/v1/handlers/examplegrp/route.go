package examplegrp

import (
	"net/http"

	authctx "github.com/OpenJenie/goserve/business/web/v1/auth"
	"github.com/OpenJenie/goserve/business/web/v1/mid"
	"github.com/OpenJenie/goserve/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth  *authctx.Auth
	Build string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build)
	authen := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, authctx.RuleAdminOnly)

	app.Handle(http.MethodGet, version, "/example", hdl.Example)
	app.Handle(http.MethodGet, version, "/exampleauth", hdl.ExampleAuth, authen, ruleAdmin)
}
