package handler

import (
	"html/template"
	"math"
	"net/url"
	"time"

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
		"dayOfWeek": func(t time.Time) string {
			return t.Format("Monday")
		},
		"uMinsToDays": func(minutes uint32) float32 {
			return float32(minutes) / 60 / 24
		},
		"minsToDays": func(minutes int32) float32 {
			return float32(roundPlaces(float64(minutes)/60/24, 2))
		},
		"shortDate": func(t time.Time) string {
			return t.Format("January 2 2006")
		},
		"urlValue": func(v url.Values, k string) string {
			return v.Get(k)
		},
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

func round(f float64) float64 {
	return math.Floor(f + .5)
}

func roundPlaces(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return round(f*shift) / shift
}
