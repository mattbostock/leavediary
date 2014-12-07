package main

import (
	_ "expvar"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
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

	n := negroni.New(negroni.NewRecovery())
	l := negronilogrus.NewMiddleware()

	n.Use(l)
	n.UseHandler(m)

	addr := ":3000"
	l.Logger.Infof("Listening on %s", addr)
	l.Logger.Fatal(http.ListenAndServe(addr, n))
}
