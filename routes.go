package main

import (
	_ "expvar"
	"net/http"

	"gitlab.com/mattbostock/timeoff/handler"
)

func registerRoutes() {
	// Expose `expvar` debug variables
	mux.Add("GET", "/debug/vars", http.DefaultServeMux)

	mux.Get("/", handler.Index)
}
