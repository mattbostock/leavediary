package main

import "gitlab.com/mattbostock/timeoff/handler"

func init() {
	m.Get("/", handler.Index)
}
