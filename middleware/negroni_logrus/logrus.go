package negroniLogrus

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

type middleware struct {
	logger *logrus.Logger
}

func New(l *logrus.Logger) *middleware {
	return &middleware{logger: l}
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	m.logger.WithFields(logrus.Fields{
		"method":  r.Method,
		"request": r.RequestURI,
		"remote":  r.RemoteAddr,
	}).Info("started handling request")

	next(w, r)

	latency := time.Since(start)
	res := w.(negroni.ResponseWriter)
	m.logger.WithFields(logrus.Fields{
		"status":      res.Status(),
		"method":      r.Method,
		"request":     r.RequestURI,
		"remote":      r.RemoteAddr,
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
	}).Info("completed handling request")
}
