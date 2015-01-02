package main

import "gitlab.com/mattbostock/timeoff/handler"

func init() {
	mux.Get("/", handler.Index)
}
