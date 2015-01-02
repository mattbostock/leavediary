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

const assetsPath = "assets"

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
	mux        = pat.New()
	n          = negroni.New()
	logHandler = negronilogrus.NewMiddleware()
)

func init() {
	// configure logging
	n.Use(logHandler)
	handler.SetLogger(logHandler.Logger)
	if config.debug {
		logHandler.Logger.Level = logrus.DebugLevel
	}

	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir(assetsPath)))
	n.Use(gzip.Gzip(gzip.BestCompression))
	n.UseHandler(mux)

	mux.Add("GET", "/debug/vars", http.DefaultServeMux)

	if config.addr == "" {
		config.addr = ":3000"
	}
}

func main() {
	logHandler.Logger.Infof("Listening on %s", config.addr)

	if config.tlsCert == "" && config.tlsKey == "" {
		logHandler.Logger.Warningln(noTLSCertificateError)
		logHandler.Logger.Fatal(http.ListenAndServe(config.addr, n))
	} else {
		logHandler.Logger.Infoln("Listening with TLS")
		logHandler.Logger.Fatal(http.ListenAndServeTLS(config.addr, config.tlsCert, config.tlsKey, n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TLS_CERT " +
	"and TLS_KEY environment variables to enable TLS."
