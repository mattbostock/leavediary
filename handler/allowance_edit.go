package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mattbostock/leavediary/model"
)

func AllowanceEdit(w http.ResponseWriter, r *http.Request) {
	var (
		allowanceID uint64
		formValues  url.Values
		userErr     string
	)
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	job := user.CurrentJob()

	if allowanceID, _ = strconv.ParseUint(r.URL.Query().Get(":id"), 10, 64); allowanceID > 0 {
		leaveAllowance, err := model.FindLeaveAllowance(allowanceID, job.ID)
		if err != nil {
			internalError(w, err)
			return
		}

		formValues = make(url.Values)
		formValues.Add("days", strconv.FormatInt((int64(leaveAllowance.Minutes)/24/60), 10))
		formValues.Add("allowance_start_day", leaveAllowance.StartTime.Format("2"))
		formValues.Add("allowance_start_month", leaveAllowance.StartTime.Format("1"))
		formValues.Add("allowance_start_year", leaveAllowance.StartTime.Format("2006"))
	}

	if r.Method == "POST" {
		formValues = r.PostForm

		days, err := strconv.ParseFloat(r.PostFormValue("days"), 32)
		if err != nil {
			internalError(w, err)
			return
		}

		allowanceStart, err := time.ParseInLocation("2006-1-2",
			fmt.Sprintf("%s-%s-%s",
				r.PostFormValue("allowance_start_year"),
				r.PostFormValue("allowance_start_month"),
				r.PostFormValue("allowance_start_day"),
			),
			user.TZLocation())

		if err != nil {
			internalError(w, err)
			return
		}

		allowanceEnd := allowanceStart.AddDate(1, 0, 0)

		a := model.LeaveAllowance{
			ID:           allowanceID,
			IsAdjustment: false,
			Description:  "Annual leave",
			StartTime:    allowanceStart,
			EndTime:      allowanceEnd,
			Minutes:      int32(days * 24 * 60),
			JobID:        user.Jobs[0].ID,
		}

		overlapping, err := a.OverlapsAnother()
		if err != nil {
			internalError(w, err)
			return
		}

		intersects, err := a.IntersectsLeaveRequest()
		if err != nil {
			internalError(w, err)
			return
		}

		// FIXME: Make error messages more useful
		switch {
		case overlapping.ID > 0:
			userErr = "The dates provided overlap with an existing leave year."
		case intersects.ID > 0:
			userErr = "The dates provided fall part-way through existing leave requests."
		default:
			if allowanceID > 0 {
				err = a.Save()
			} else {
				user.Jobs[0].LeaveAllowances = append(user.Jobs[0].LeaveAllowances, a)
				err = user.Save()

			}

			if err != nil {
				internalError(w, err)
				return
			}

			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}

	output.HTML(w, http.StatusOK, "allowance_edit",
		&struct {
			ID         uint64
			FormValues url.Values
			User       model.User
			UserErr    string
		}{
			allowanceID,
			formValues,
			user,
			userErr,
		})
}
