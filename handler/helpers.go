package handler

import (
	"net/http"
	"github.com/nbio/httpcontext"
	"gitlab.com/mattbostock/timeoff/model"
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
