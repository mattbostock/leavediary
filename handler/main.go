package handler

import (
	"html/template"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/mattbostock/timeoff/middleware/sessions"
	"github.com/unrolled/render"
	"golang.org/x/oauth2"
)

const oauthStateCookieName = "github_state"

var (
	githubAPIBaseURL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/",
	}
	log            *logrus.Logger
	oauthConfig    *oauth2.Config
	output         *render.Render
	sessionManager *sessions.Manager
	version        string
)

func init() {
	templateFuncs := &template.FuncMap{
		"version": func() string {
			return version
		},
	}

	output = render.New(render.Options{
		Funcs:  []template.FuncMap{*templateFuncs},
		Layout: "layout",
	})
}

func SetLogger(l *logrus.Logger) {
	log = l
}

func SetOauthConfig(o *oauth2.Config) {
	oauthConfig = o
}

func SetSessionManager(s *sessions.Manager) {
	sessionManager = s
}

func SetVersion(v string) {
	version = v
}
