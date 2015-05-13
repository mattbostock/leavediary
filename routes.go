package main

import (
	_ "expvar"
	"net/http"

	"github.com/mattbostock/leavediary/handler"
)

func registerRoutes() {
	// Expose `expvar` debug variables
	mux.Get("/debug/vars", http.DefaultServeMux)
	mux.Get("/allowance", http.HandlerFunc(handler.Allowance))
	mux.Get("/allowance/delete/:id", http.HandlerFunc(handler.AllowanceDelete))
	mux.Get("/allowance/new", http.HandlerFunc(handler.AllowanceEdit))
	mux.Get("/allowance/:id", http.HandlerFunc(handler.AllowanceEdit))

	mux.Post("/allowance/new", http.HandlerFunc(handler.AllowanceEdit))
	mux.Post("/allowance/:id", http.HandlerFunc(handler.AllowanceEdit))

	mux.Get("/request/delete/:id", http.HandlerFunc(handler.RequestDelete))
	mux.Get("/request/new", http.HandlerFunc(handler.RequestEdit))
	mux.Get("/request/:id", http.HandlerFunc(handler.RequestEdit))

	mux.Post("/request/new", http.HandlerFunc(handler.RequestEdit))
	mux.Post("/request/:id", http.HandlerFunc(handler.RequestEdit))

	mux.Get("/settings", http.HandlerFunc(handler.Settings))
	mux.Post("/settings", http.HandlerFunc(handler.Settings))

	mux.Get("/export/ics/:secret", http.HandlerFunc(handler.ExportICS))
	mux.Get("/export/csv", http.HandlerFunc(handler.ExportCSV))

	mux.Get("/oauth/github/callback", http.HandlerFunc(handler.GithubOauthCallback))
	mux.Get("/dashboard", http.HandlerFunc(handler.Dashboard))
	mux.Get("/logout", http.HandlerFunc(handler.Logout))

	mux.Get("/", http.HandlerFunc(handler.Index))
}
