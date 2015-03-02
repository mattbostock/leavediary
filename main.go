package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"runtime"
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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	assetsPath  = "assets"
	defaultAddr = "localhost:3000"
	sessionName = "timeoff_session"
)

var (
	config = &struct {
		addr               string
		allowedHosts       []string
		cookieHashKey      []byte
		debug              bool
		gitHubClientID     string
		gitHubClientSecret string
		tlsCert            string
		tlsKey             string
	}{
		addr:               os.Getenv("ADDR"),
		cookieHashKey:      []byte(os.Getenv("COOKIE_KEY")),
		debug:              os.Getenv("DEBUG") != "",
		gitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		gitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		tlsCert:            os.Getenv("TLS_CERT"),
		tlsKey:             os.Getenv("TLS_KEY"),
	}

	mux = pat.New()
	log = logrus.New()
	version = ""
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if config.debug {
		log.Level = logrus.DebugLevel
	}

	if config.addr == "" {
		config.addr = defaultAddr
	}
}

func main() {
	if version == "" {
		log.Fatalln(errMakeFileNotUsed)
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

	if config.gitHubClientID == "" || config.gitHubClientSecret == "" {
		log.Fatalf(errNoGitHubCredentials)
	}

	handler.SetOauthConfig(&oauth2.Config{
		ClientID:     config.gitHubClientID,
		ClientSecret: config.gitHubClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	})

	model.SetLogger(log)
	model.InitDB("sqlite3", "sqlite.db")

	sessions.SetLogger(log)
	handler.SetLogger(log)

	sessionManager := sessions.New(sessionName, config.cookieHashKey)
	handler.SetSessionManager(sessionManager)

	secureOpts := secure.Options{
		AllowedHosts:          config.allowedHosts,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; img-src 'self' data:",
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

	c := &tls.Config{MinVersion: tls.VersionTLS10} // disable SSLv3, prevent POODLE attack
	s := &http.Server{Addr: config.addr, Handler: n, TLSConfig: c}

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

	errCookieHashKeyWrongLength = "COOKIE_KEY environment variable must be 32 characters long. Length provided: %d"

	errMakeFileNotUsed = "Makefile was not used when compiling binary, run 'make' to re-compile"

	errNoCookieHashKey = "No cookie hash key supplied. You should set the COOKIE_KEY " +
		"environment variable in a production environment. Falling back to use a temporary key " +
		"which will persist only for the current running process."

	errNoGitHubCredentials = "No GitHub Oauth credentials supplied. Set both GITHUB_CLIENT_ID and " +
		"GITHUB_CLIENT_SECRET environment variables."

	errNoTLSCertificate = "No TLS certficiate supplied. Consider setting TLS_CERT " +
		"and TLS_KEY environment variables to enable TLS. TimeOff will not work unless you " +
		"are using TLS upstream."
)
