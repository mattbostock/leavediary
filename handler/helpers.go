package handler

import (
	"net/http"

	"github.com/mattbostock/leavediary/model"
	"github.com/nbio/httpcontext"
)

func internalError(w http.ResponseWriter, err error) {
	log.Error(err)
	output.HTML(w, http.StatusInternalServerError, "error",
		&struct {
			Msg  string
			User model.User
		}{
			"",
			model.User{},
		})
}

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
