package main

import (
	_ "expvar"
	"net/http"

	"github.com/mattbostock/timeoff/handler"
)

func registerRoutes() {
	// Expose `expvar` debug variables
	mux.Get("/debug/vars", http.DefaultServeMux)

	mux.Get("/dashboard/settings", http.HandlerFunc(handler.DashboardSettings))
	mux.Post("/dashboard/settings", http.HandlerFunc(handler.DashboardSettings))

	mux.Get("/oauth/github/callback", http.HandlerFunc(handler.GithubOauthCallback))
	mux.Get("/dashboard", http.HandlerFunc(handler.Dashboard))
	mux.Get("/logout", http.HandlerFunc(handler.Logout))

	mux.Get("/", http.HandlerFunc(handler.Index))
}
