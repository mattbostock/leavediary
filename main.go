package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/bradfitz/http2"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"gitlab.com/mattbostock/timeoff/handler"
	"gitlab.com/mattbostock/timeoff/middleware/negroni_logrus"
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

	db  = initDB()
	mux = pat.New()
	log = logrus.New()
)

func init() {
	if config.debug {
		log.Level = logrus.DebugLevel
	}

	if config.addr == "" {
		config.addr = defaultAddr
	}

	handler.SetLogger(log)

}

func main() {
	n := negroni.New()
	n.Use(negroniLogrus.New(log)) // logger must be first middleware
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir(assetsPath)))
	n.UseHandler(mux)
	registerRoutes()

	log.Infof("Listening on %s", config.addr)

	s := &http.Server{Addr: config.addr, Handler: n}

	if config.tlsCert == "" && config.tlsKey == "" {
		log.Warningln(noTLSCertificateError)
		log.Fatal(s.ListenAndServe())
	} else {
		http2.ConfigureServer(s, nil)
		log.Infoln("TLS-only; HTTP/2 enabled")
		log.Fatal(s.ListenAndServeTLS(config.tlsCert, config.tlsKey))
	}
}

const noTLSCertificateError = "No TLS certficiate supplied. Consider setting TLS_CERT " +
	"and TLS_KEY environment variables to enable TLS."
