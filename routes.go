package main

import (
	_ "expvar"
	"net/http"

	"gitlab.com/mattbostock/timeoff/handler"
)

func registerRoutes() {
	// Expose `expvar` debug variables
	mux.Handle("/debug/vars", http.DefaultServeMux)

	mux.Get("/dashboard", handler.Dashboard)
	mux.Get("/", handler.Index)
}
