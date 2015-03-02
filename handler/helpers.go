package handler

import (
	"github.com/nbio/httpcontext"
	"gitlab.com/mattbostock/timeoff/model"
	"net/http"
)

func showError(w http.ResponseWriter, msg string, statusCode int) {
	output.HTML(w, statusCode, "error", msg)
}

func currentUser(r *http.Request) model.User {
	user := httpcontext.Get(r, "user")

	switch user := user.(type) {
	case model.User:
		return user
	default:
		return model.User{}
	}
}
