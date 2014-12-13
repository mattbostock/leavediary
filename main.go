package main

import (
	_ "expvar"
	"github.com/codegangsta/negroni"
	"github.com/meatballhat/negroni-logrus"
	"github.com/unrolled/render"
	"net/http"
	"os"
)

func main() {
	Run()
}

func Run() {
	m := http.DefaultServeMux
	n := negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l := negronilogrus.NewMiddleware()
	r := render.New(render.Options{
		Layout: "layout",
	})

	n.Use(l)
	n.UseHandler(m)

	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", "world")
	})

	var addr string
	if len(os.Getenv("PORT")) > 0 {
		addr = ":" + os.Getenv("PORT")
	} else {
		addr = ":3000"
	}

	l.Logger.Infof("Listening on %s", addr)
	l.Logger.Fatal(http.ListenAndServe(addr, n))
}
