package handler

import (
	"github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
	"gitlab.com/mattbostock/timeoff/middleware/sessions"
	"golang.org/x/oauth2"
)

const oauthStateCookieName = "github_state"

var (
	log            *logrus.Logger
	oauthConfig    *oauth2.Config
	output         = render.New(render.Options{Layout: "layout"})
	sessionManager *sessions.Manager
)

func SetLogger(l *logrus.Logger) {
	log = l
}

func SetOauthConfig(o *oauth2.Config) {
	oauthConfig = o
}

func SetSessionManager(s *sessions.Manager) {
	sessionManager = s
}
