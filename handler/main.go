package handler

import (
	"github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
	"gitlab.com/mattbostock/timeoff/middleware/sessions"
)

const oauthStateCookieName = "github_state"

var (
	log            *logrus.Logger
	output         = render.New(render.Options{Layout: "layout"})
	sessionManager *sessions.Manager
)

func SetLogger(l *logrus.Logger) {
	log = l
}

func SetSessionManager(s *sessions.Manager) {
	sessionManager = s
}
