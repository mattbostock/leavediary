package handler

import (
	"net/http"

	"github.com/mattbostock/timeoff/model"
	"github.com/nbio/httpcontext"
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
