package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/bradfitz/http2"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	"github.com/gorilla/securecookie"
	"github.com/unrolled/secure"
	"gitlab.com/mattbostock/timeoff/handler"
	"gitlab.com/mattbostock/timeoff/middleware/negroni_logrus"
	"gitlab.com/mattbostock/timeoff/middleware/sessions"
	"gitlab.com/mattbostock/timeoff/model"
)

const (
	assetsPath  = "assets"
	defaultAddr = ":3000"
	sessionName = "timeoff_session"
)

var (
	config = &struct {
		addr               string
		allowedHosts       []string
		cookieHashKey      []byte
		debug              bool
		tlsCert            string
		tlsKey             string
	}{
		addr:               os.Getenv("ADDR"),
		cookieHashKey:      []byte(os.Getenv("COOKIE_KEY")),
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

	if len(config.cookieHashKey) == 0 {
		log.Warningln(errNoCookieHashKey)
		config.cookieHashKey = securecookie.GenerateRandomKey(32)
	}
	if len(config.cookieHashKey) != 32 {
		// additonal check as securecookie.GenerateRandomKey() does not return errors
		log.Fatalf(errCookieHashKeyWrongLength, len(config.cookieHashKey))
	}

	model.SetLogger(log)
	model.InitDB("sqlite3", "sqlite.db")

	sessions.SetLogger(log)
	handler.SetLogger(log)

}

func main() {
	sessionManager := sessions.New(sessionName, config.cookieHashKey)
	handler.SetSessionManager(sessionManager)

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
	n.UseHandler(sessionManager)
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

	errCookieHashKeyWrongLength = "COOKIE_KEY environment variable must be 32 characters long. Length provided: %d"

	errNoCookieHashKey = "No cookie hash key supplied. You should set the COOKIE_KEY " +
		"environment variable in a production environment. Falling back to use a temporary key " +
		"which will persist only for the current running process."
)
