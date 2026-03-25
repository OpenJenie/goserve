package checkgrp

import (
	"net/http"

	"github.com/OpenJenie/goserve/foundation/logger"
	"github.com/OpenJenie/goserve/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	hdl := New(cfg.Build, cfg.Log)
	app.HandleNoMiddleware(http.MethodGet, version, "/readiness", hdl.Readiness)
	app.HandleNoMiddleware(http.MethodGet, version, "/liveness", hdl.Liveness)
}
