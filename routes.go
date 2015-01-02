package main

import "gitlab.com/mattbostock/timeoff/handler"

func registerRoutes() {
	mux.Get("/", handler.Index)
}
