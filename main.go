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
	config = &struct {
		addr    string
		debug   bool
		tlsCert string
		tlsKey  string
	}{
		addr:    os.Getenv("ADDR"),
		debug:   len(os.Getenv("DEBUG")) > 0,
		tlsCert: os.Getenv("TLS_CERT"),
		tlsKey:  os.Getenv("TLS_KEY"),
	}
	m = pat.New()
	n = negroni.New(negroni.NewRecovery(), negroni.NewStatic(http.Dir("assets")))
	l = negronilogrus.NewMiddleware()
)

func init() {
	// configure logging
	n.Use(l)
	handler.SetLogger(l.Logger)
	if config.debug {
		l.Logger.Level = logrus.DebugLevel
	}

	n.Use(gzip.Gzip(gzip.BestCompression))
	n.UseHandler(m)

	m.Add("GET", "/debug/vars", http.DefaultServeMux)

	if config.addr == "" {
		config.addr = ":3000"
	}
}

func main() {
	l.Logger.Infof("Listening on %s", config.addr)

	if config.tlsCert == "" && config.tlsKey == "" {
		l.Logger.Warningln(noTLSCertificateError)
		l.Logger.Fatal(http.ListenAndServe(config.addr, n))
	} else {
		l.Logger.Infoln("Listening with TLS")
		l.Logger.Fatal(http.ListenAndServeTLS(config.addr, config.tlsCert, config.tlsKey, n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TLS_CERT " +
	"and TLS_KEY environment variables to enable TLS."
