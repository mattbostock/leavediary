package main

import (
	_ "expvar"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/meatballhat/negroni-logrus"
)

var (
	addr = os.Getenv("TIMEOFF_ADDR")
	m    = pat.New()
	n    = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l    = negronilogrus.NewMiddleware()
)

func init() {
	n.Use(l)
	n.UseHandler(m)

	m.Add("GET", "/debug/vars", http.DefaultServeMux)

	if addr == "" {
		addr = ":3000"
	}
}

func main() {
	tlsCert := os.Getenv("TIMEOFF_TLS_CERT")
	tlsKey := os.Getenv("TIMEOFF_TLS_KEY")

	l.Logger.Infof("Listening on %s", addr)

	if tlsCert == "" && tlsKey == "" {
		l.Logger.Warningln(noTLSCertificateError)
		l.Logger.Fatal(http.ListenAndServe(addr, n))
	} else {
		l.Logger.Infoln("Listening with TLS")
		l.Logger.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TIMEOFF_TLS_CERT " +
	"and TIMEOFF_TLS_KEY environment variables to enable TLS."
