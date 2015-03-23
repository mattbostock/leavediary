package handler

import (
	"net/http"

	"github.com/mattbostock/timeoff/model"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	output.HTML(w, http.StatusOK, "dashboard", &struct {
		User        model.User
		Employments []model.Employment
	}{
		user,
		nil,
	})
}
