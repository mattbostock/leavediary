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

var (
	addr = os.Getenv("TIMEOFF_ADDR")
	m    = pat.New()
	n    = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l    = negronilogrus.NewMiddleware()
       o    = render.New(render.Options{
                      Layout: "layout",
              })
)

func init() {
	n.Use(l)
	n.UseHandler(m)

	m.Get("/debug/vars", http.DefaultServeMux)

	m.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.HTML(w, http.StatusOK, "index", "world")
	}))

	if addr == "" {
		addr = ":3000"
	}
}

func main() {
	l.Logger.Infof("Listening on %s", addr)
	if os.Getenv("TIMEOFF_TLS_CERT") == "" && os.Getenv("TIMEOFF_TLS_KEY") == "" {
		l.Logger.Warningln(noTLSCertificateError)
		l.Logger.Fatal(http.ListenAndServe(addr, n))
	} else {
		l.Logger.Infoln("Listening with TLS")
		l.Logger.Fatal(http.ListenAndServeTLS(addr, os.Getenv("TIMEOFF_TLS_CERT"), os.Getenv("TIMEOFF_TLS_KEY"), n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TIMEOFF_TLS_CERT " +
	"and TIMEOFF_TLS_KEY environment variables to enable TLS."
