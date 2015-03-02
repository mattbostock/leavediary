package handler

import (
	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
	"gitlab.com/mattbostock/timeoff/middleware/sessions"
	"golang.org/x/oauth2"
)

const oauthStateCookieName = "github_state"

var (
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
