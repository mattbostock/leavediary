package handler

import (
	"net/http"

	"github.com/mattbostock/leavediary/model"
)

func Allowance(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	job := user.CurrentJob()

	leavePeriods, err := job.LeavePeriods()
	if err != nil {
		internalError(w, err)
		return
	}

	output.HTML(w, http.StatusOK, "allowance",
		&struct {
			JobID        uint64
			LeavePeriods []model.LeaveAllowance
			User         model.User
		}{
			job.ID,
			leavePeriods,
			user,
		})
}
