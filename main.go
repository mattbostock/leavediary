package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/meatballhat/negroni-logrus"
	"github.com/phyber/negroni-gzip/gzip"
	"gitlab.com/mattbostock/timeoff/handler"
)

const (
	assetsPath  = "assets"
	defaultAddr = ":3000"
)

var (
	config = &struct {
		addr    string
		debug   bool
		tlsCert string
		tlsKey  string
	}{
		addr:    os.Getenv("ADDR"),
		debug:   os.Getenv("DEBUG") != "",
		tlsCert: os.Getenv("TLS_CERT"),
		tlsKey:  os.Getenv("TLS_KEY"),
	}

	mux        = pat.New()
	logHandler = negronilogrus.NewMiddleware()
	log        = logHandler.Logger
)

func init() {
	if config.addr == "" {
		config.addr = defaultAddr
	}

	if config.debug {
		log.Level = logrus.DebugLevel
	}

	handler.SetLogger(log)
}

func main() {
	n := negroni.New()
	n.Use(logHandler) // logger must be first middleware
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir(assetsPath)))
	n.Use(gzip.Gzip(gzip.BestCompression))

	n.UseHandler(mux)
	registerRoutes()

	log.Infof("Listening on %s", config.addr)

	if config.tlsCert == "" && config.tlsKey == "" {
		log.Warningln(noTLSCertificateError)
		log.Fatal(http.ListenAndServe(config.addr, n))
	} else {
		log.Infoln("Listening with TLS")
		log.Fatal(http.ListenAndServeTLS(config.addr, config.tlsCert, config.tlsKey, n))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TLS_CERT " +
	"and TLS_KEY environment variables to enable TLS."
