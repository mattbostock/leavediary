package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bradfitz/http2"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/unrolled/secure"
	"gitlab.com/mattbostock/timeoff/handler"
	"gitlab.com/mattbostock/timeoff/middleware/negroni_logrus"
	"gitlab.com/mattbostock/timeoff/model"
)

const (
	assetsPath  = "assets"
	defaultAddr = ":3000"
)

var (
	config = &struct {
		addr               string
		allowedHosts       []string
		debug              bool
		tlsCert            string
		tlsKey             string
	}{
		addr:               os.Getenv("ADDR"),
		debug:              os.Getenv("DEBUG") != "",
		tlsCert:            os.Getenv("TLS_CERT"),
		tlsKey:             os.Getenv("TLS_KEY"),
	}

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

	if os.Getenv("ALLOWED_HOSTS") == "" {
		config.allowedHosts = nil
	} else {
		config.allowedHosts = strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")
	}

	model.SetLogger(log)
	model.InitDB("sqlite3", "sqlite.db")

	handler.SetLogger(log)

}

func main() {
	secureOpts := secure.Options{
		AllowedHosts:          config.allowedHosts,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline'; img-src 'self' data:",
		FrameDeny:             true,
		STSIncludeSubdomains:  true,
		STSSeconds:            365 * 24 * 60 * 60,
	}
	secureMiddleware := secure.New(secureOpts)

	n := negroni.New()
	n.Use(negroniLogrus.New(log)) // logger must be first middleware
	n.Use(negroni.NewRecovery())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(negroni.NewStatic(http.Dir(assetsPath)))
	n.UseHandler(mux)
	registerRoutes()

	log.Infof("Listening on %s", config.addr)

	if config.allowedHosts != nil {
		log.Infof("Allowed HTTP hosts: %q", config.allowedHosts)
	} else {
		log.Warningln(errAllHostsEnabled)
	}

	s := &http.Server{Addr: config.addr, Handler: n}

	if config.tlsCert == "" && config.tlsKey == "" {
		log.Warningln(errNoTLSCertificate)
		log.Fatal(s.ListenAndServe())
	} else {
		http2.ConfigureServer(s, nil)
		log.Infoln("TLS-only; HTTP/2 enabled")
		log.Fatal(s.ListenAndServeTLS(config.tlsCert, config.tlsKey))
	}
}

const (
	errAllHostsEnabled = "Accepting connections for all HTTP hosts. Consider setting the ALLOWED_HOSTS environment variable."

       errNoTLSCertificate = "No TLS certficiate supplied. Consider setting TLS_CERT " +
		"and TLS_KEY environment variables to enable TLS."
)
