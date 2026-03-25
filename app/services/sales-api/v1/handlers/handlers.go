package handlers

import (
	"github.com/OpenJenie/goserve/app/services/sales-api/v1/handlers/checkgrp"
	"github.com/OpenJenie/goserve/app/services/sales-api/v1/handlers/examplegrp"
	v1 "github.com/OpenJenie/goserve/business/web/v1"
	"github.com/OpenJenie/goserve/foundation/web"
)

type Routes struct{}

// Add implements the RouterAdder interface.
func (Routes) Add(app *web.App, apiCfg v1.APIMuxConfig) {
	examplegrp.Routes(app, examplegrp.Config{
		Auth:  apiCfg.Auth,
		Build: apiCfg.Build,
	})

	checkgrp.Routes(app, checkgrp.Config{
		Build: apiCfg.Build,
		Log:   apiCfg.Log,
	})
}
