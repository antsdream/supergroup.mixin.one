package main

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/MixinMessenger/supergroup.mixin.one/config"
	"github.com/MixinMessenger/supergroup.mixin.one/durable"
	"github.com/MixinMessenger/supergroup.mixin.one/middlewares"
	"github.com/MixinMessenger/supergroup.mixin.one/routes"
	"github.com/dimfeld/httptreemux"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	"github.com/unrolled/render"
)

func StartServer(spanner *spanner.Client) error {
	logger := durable.NewLoggerClient()
	router := httptreemux.New()
	routes.RegisterHanders(router)
	routes.RegisterRoutes(router)
	handler := middlewares.Authenticate(router)
	handler = middlewares.Constraint(handler)
	handler = middlewares.Context(handler, spanner, render.New(render.Options{UnEscapeHTML: true}))
	handler = middlewares.Stats(handler, "http", true, config.BuildVersion)
	handler = middlewares.Log(handler, logger, "http")
	handler = handlers.ProxyHeaders(handler)

	return gracehttp.Serve(&http.Server{Addr: fmt.Sprintf(":%d", config.HTTPListenPort), Handler: handler})
}
