package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/Sirupsen/logrus"
	"github.com/bmizerany/pat"
	"github.com/bradfitz/http2"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/securecookie"
	"github.com/mattbostock/leavediary/handler"
	"github.com/mattbostock/leavediary/middleware/negroni_logrus"
	"github.com/mattbostock/leavediary/middleware/sessions"
	"github.com/mattbostock/leavediary/model"
	"github.com/unrolled/secure"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	assetsPath  = "assets"
	defaultAddr = "localhost:3000"
	sessionName = "leavediary_session"
)

var (
	config = &struct {
		addr               string
		allowedHosts       []string
		cookieHashKey      []byte
		debug              bool
		dbDialect          string
		dbDataSource       string
		gitHubClientID     string
		gitHubClientSecret string
		rateLimitPerMin    uint8
		tlsCert            string
		tlsKey             string
	}{
		addr:               os.Getenv("ADDR"),
		cookieHashKey:      []byte(os.Getenv("COOKIE_KEY")),
		dbDialect:          os.Getenv("DB_DIALECT"),
		dbDataSource:       os.Getenv("DB_DATASOURCE"),
		debug:              os.Getenv("DEBUG") != "",
		gitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		gitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		tlsCert:            os.Getenv("TLS_CERT"),
		tlsKey:             os.Getenv("TLS_KEY"),
	}

	mux     = pat.New()
	log     = logrus.New()
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

	rate, _ := strconv.ParseUint(os.Getenv("RATE_LIMIT_PER_MIN"), 10, 8)
	config.rateLimitPerMin = uint8(rate)

	if rate == 0 {
		config.rateLimitPerMin = 240
	}
}

func main() {
	if version == "" {
		log.Fatalln(errMakeFileNotUsed)
	}

	v := flag.Bool("version", false, "prints current version")
	flag.Parse()

	if *v {
		fmt.Println(version)
		os.Exit(0)
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

	if config.dbDialect == "" && config.dbDataSource == "" {
		config.dbDialect = "sqlite3"
		config.dbDataSource = ":memory:"
	}

	handler.SetVersion(version)

	handler.SetOauthConfig(&oauth2.Config{
		ClientID:     config.gitHubClientID,
		ClientSecret: config.gitHubClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	})

	model.SetLogger(log)
	model.InitDB(config.dbDialect, config.dbDataSource)

	sessions.SetLogger(log)
	handler.SetLogger(log)

	sessionManager := sessions.New(sessionName, config.cookieHashKey)
	handler.SetSessionManager(sessionManager)

	secureOpts := secure.Options{
		AllowedHosts:          config.allowedHosts,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'sha256-BWV1eSks2QM8blQZAbrSRSwqg3VFfmJ2d6r7yBVBXGY='; style-src 'self' 'unsafe-inline'; img-src 'self' data:",
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

	// throttle requests by remote IP though X-Forwarded-For could be spoofed
	varyHost := func(r *http.Request) string {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		return host + r.Header.Get("X-Forwarded-For")
	}

	t := throttled.RateLimit(throttled.PerMin(config.rateLimitPerMin), &throttled.VaryBy{Custom: varyHost}, store.NewMemStore(1000))
	t.DeniedHandler = http.HandlerFunc(handler.TooManyRequests)
	log.Infof("Throttling requests at %d per minute per remote IP address", config.rateLimitPerMin)
	h := t.Throttle(n)

	log.Infof("Listening on %s", config.addr)

	if config.allowedHosts != nil {
		log.Infof("Allowed HTTP hosts: %q", config.allowedHosts)
	} else {
		log.Warningln(errAllHostsEnabled)
	}

	c := &tls.Config{MinVersion: tls.VersionTLS10} // disable SSLv3, prevent POODLE attack
	s := &http.Server{Addr: config.addr, Handler: h, TLSConfig: c}

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
		"and TLS_KEY environment variables to enable TLS. LeaveDiary will not work unless you " +
		"are using TLS upstream."
)
