package main

import (
	_ "expvar"
	"net/http"
	"os"

	"github.com/bmizerany/pat"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"github.com/unrolled/render"
)

func main() {
	m := pat.New()
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l := negronilogrus.NewMiddleware()
	o := render.New(render.Options{
		Layout: "layout",
	})

	n.Use(l)
	n.UseHandler(m)

	m.Get("/debug/vars", http.DefaultServeMux)

	m.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.HTML(w, http.StatusOK, "index", "world")
	}))

	var addr string

	if os.Getenv("TIMEOFF_ADDR") == "" {
		addr = ":3000"
	} else {
		addr = os.Getenv("TIMEOFF_ADDR")
	}

	l.Logger.Infof("Listening on %s", addr)
	l.Logger.Fatal(http.ListenAndServe(addr, n))
}
