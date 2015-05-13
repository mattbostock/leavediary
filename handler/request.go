package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mattbostock/leavediary/model"
)

func Request(w http.ResponseWriter, r *http.Request) {
	var (
		requestID  uint64
		formValues url.Values
		userErr    string
	)
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if requestID, _ = strconv.ParseUint(r.URL.Query().Get(":id"), 10, 64); requestID > 0 {
		leaveRequest, err := model.FindLeaveRequest(requestID)
		if err != nil {
			internalError(w, err)
			return
		}

		formValues = make(url.Values)
		formValues.Add("days", strconv.FormatUint((uint64(leaveRequest.Minutes)/24/60), 10))
		formValues.Add("description", leaveRequest.Description)
		formValues.Add("leave_end_day", leaveRequest.EndTime.Format("2"))
		formValues.Add("leave_end_month", leaveRequest.EndTime.Format("1"))
		formValues.Add("leave_end_year", leaveRequest.EndTime.Format("2006"))
		formValues.Add("leave_start_day", leaveRequest.StartTime.Format("2"))
		formValues.Add("leave_start_month", leaveRequest.StartTime.Format("1"))
		formValues.Add("leave_start_year", leaveRequest.StartTime.Format("2006"))
	}

	if r.Method == "POST" {
		formValues = r.PostForm

		days, err := strconv.ParseFloat(r.PostFormValue("days"), 32)
		if err != nil {
			internalError(w, err)
			return
		}

		leaveStart, err := time.ParseInLocation("2006-1-2",
			fmt.Sprintf("%s-%s-%s",
				r.PostFormValue("leave_start_year"),
				r.PostFormValue("leave_start_month"),
				r.PostFormValue("leave_start_day"),
			),
			user.TZLocation())

		if err != nil {
			internalError(w, err)
			return
		}

		leaveEnd, err := time.ParseInLocation("2006-1-2",
			fmt.Sprintf("%s-%s-%s",
				r.PostFormValue("leave_end_year"),
				r.PostFormValue("leave_end_month"),
				r.PostFormValue("leave_end_day"),
			),
			user.TZLocation())

		if err != nil {
			internalError(w, err)
			return
		}

		// FIXME check if working days exceeds difference between start and end dates
		// FIXME check if leave requested starts before job starts or ends after job end date

		l := model.LeaveRequest{
			Description: r.PostFormValue("description"),
			StartTime:   leaveStart,
			EndTime:     leaveEnd,
			Minutes:     uint32(days * 24 * 60),
			JobID:       user.Jobs[0].ID,
		}

		matching, err := l.FitsExistingAllowancePeriod()
		if err != nil {
			internalError(w, err)
			return
		}

		// FIXME show more helpful message - close to current leave period or spans two?
		if matching == 0 {
			userErr = `The dates you requested are outside the leave periods
				currently defined.`
		} else {
			if requestID > 0 {
				l.ID = requestID
				l.Save()
			} else {
				user.Jobs[0].LeaveRequests = append(user.Jobs[0].LeaveRequests, l)
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

	output.HTML(w, http.StatusOK, "request",
		&struct {
			ID         uint64
			FormValues url.Values
			User       model.User
			UserErr    string
		}{
			requestID,
			formValues,
			user,
			userErr,
		})
}
