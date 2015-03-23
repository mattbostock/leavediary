package handler

import (
	"net/http"

	"github.com/mattbostock/timeoff/model"
	"github.com/nbio/httpcontext"
)

func showError(w http.ResponseWriter, statusCode int, msg string) {
	output.HTML(w, statusCode, "error",
		&struct {
			Msg  string
			User model.User
		}{
			msg,
			model.User{},
		})
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
