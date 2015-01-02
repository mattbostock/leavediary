package main

import (
	_ "expvar"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/meatballhat/negroni-logrus"
	"github.com/phyber/negroni-gzip/gzip"
	"gitlab.com/mattbostock/timeoff/handler"
)

var (
	addr = os.Getenv("ADDR")
	m    = pat.New()
	n    = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l    = negronilogrus.NewMiddleware()
)

func init() {
	n.Use(gzip.Gzip(gzip.BestCompression))
	n.Use(l)
	n.UseHandler(m)
	handler.SetLogger(l.Logger)

	if len(os.Getenv("DEBUG")) > 0 {
		l.Logger.Level = logrus.DebugLevel
	}

	m.Add("GET", "/debug/vars", http.DefaultServeMux)

	if addr == "" {
		addr = ":3000"
	}
}

func main() {
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	l.Logger.Infof("Listening on %s", addr)

	if tlsCert == "" && tlsKey == "" {
		l.Logger.Warningln(noTLSCertificateError)
		l.Logger.Fatal(http.ListenAndServe(addr, n))
	} else {
		l.Logger.Infoln("Listening with TLS")
		l.Logger.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TLS_CERT " +
	"and TLS_KEY environment variables to enable TLS."
