package main

import (
	_ "expvar"
	"fmt"
	"github.com/codegangsta/negroni"
	"net/http"
)

func main() {
	Run()
}

func Run() {
	m := http.DefaultServeMux

	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	n := negroni.Classic()
	n.UseHandler(m)
	n.Run(":3000")
}
