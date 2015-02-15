package main

import (
	_ "expvar"
	"net/http"

	"gitlab.com/mattbostock/timeoff/handler"
)

func registerRoutes() {
	mux.NotFoundHandler = http.NotFoundHandler()

	// Expose `expvar` debug variables
	mux.Handle("/debug/vars", http.DefaultServeMux)

	mux.Get("/oauth/github/callback", handler.GithubOauthCallback)
	mux.Get("/dashboard", handler.Dashboard)
	mux.Get("/logout", handler.Logout)

	mux.Handle("/", http.HandlerFunc(handler.Index))
}
