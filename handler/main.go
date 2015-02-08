package handler

import (
	"github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
)

var (
	log    *logrus.Logger
	output = render.New(render.Options{Layout: "layout"})
)

func SetLogger(l *logrus.Logger) {
	log = l
}
