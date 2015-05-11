package main

import (
	_ "expvar"
	"net/http"

	"github.com/mattbostock/timeoff/handler"
)

func registerRoutes() {
	// Expose `expvar` debug variables
	mux.Get("/debug/vars", http.DefaultServeMux)
	mux.Get("/dashboard/allowance/delete/:id", http.HandlerFunc(handler.DashboardAllowanceDelete))
	mux.Get("/dashboard/allowance/new", http.HandlerFunc(handler.DashboardAllowance))
	mux.Get("/dashboard/allowance/:id", http.HandlerFunc(handler.DashboardAllowance))

	mux.Post("/dashboard/allowance/new", http.HandlerFunc(handler.DashboardAllowance))
	mux.Post("/dashboard/allowance/:id", http.HandlerFunc(handler.DashboardAllowance))

	mux.Get("/dashboard/request/delete/:id", http.HandlerFunc(handler.DashboardRequestDelete))
	mux.Get("/dashboard/request/new", http.HandlerFunc(handler.DashboardRequest))
	mux.Get("/dashboard/request/:id", http.HandlerFunc(handler.DashboardRequest))

	mux.Post("/dashboard/request/new", http.HandlerFunc(handler.DashboardRequest))
	mux.Post("/dashboard/request/:id", http.HandlerFunc(handler.DashboardRequest))

	mux.Get("/dashboard/settings", http.HandlerFunc(handler.DashboardSettings))
	mux.Post("/dashboard/settings", http.HandlerFunc(handler.DashboardSettings))

	mux.Get("/export/ics/:secret", http.HandlerFunc(handler.ExportICS))
	mux.Get("/export/csv", http.HandlerFunc(handler.ExportCSV))

	mux.Get("/oauth/github/callback", http.HandlerFunc(handler.GithubOauthCallback))
	mux.Get("/dashboard", http.HandlerFunc(handler.Dashboard))
	mux.Get("/logout", http.HandlerFunc(handler.Logout))

	mux.Get("/", http.HandlerFunc(handler.Index))
}
