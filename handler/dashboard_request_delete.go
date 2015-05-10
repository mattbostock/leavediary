package handler

import (
	"net/http"
	"strconv"

	"github.com/mattbostock/timeoff/model"
)

func DashboardRequestDelete(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	id, _ := strconv.ParseUint(r.URL.Query().Get(":id"), 10, 64)
	err := model.DeleteLeaveRequest(id)
	if err != nil {
		internalError(w, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}
