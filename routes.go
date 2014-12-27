package main

import (
	"net/http"

	"gitlab.com/mattbostock/timeoff/handler"
)

func init() {
	m.Get("/", http.HandlerFunc(handler.Index))
}
