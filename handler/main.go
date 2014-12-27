package handler

import (
	"github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
)

var (
	l *logrus.Logger
	o = render.New(render.Options{
		Layout: "layout",
	})
)

func SetLogger(logger *logrus.Logger) {
	l = logger
}
