package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/mattbostock/leavediary/model"
)

func Settings(w http.ResponseWriter, r *http.Request) {
	var formValues url.Values
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	job := user.CurrentJob()
	if job.ID > 0 {
		formValues = make(url.Values)
		formValues.Add("employer_name", job.EmployerName)
		formValues.Add("job_start_year", job.StartTime.Format("2006"))
		formValues.Add("job_start_month", job.StartTime.Format("1"))
		formValues.Add("job_start_day", job.StartTime.Format("2"))
	}

	if r.Method == "POST" {
		formValues = r.PostForm

		employerName := r.PostFormValue("employer_name")
		jobStart, err := time.ParseInLocation("2006-1-2",
			fmt.Sprintf("%s-%s-%s",
				r.PostFormValue("job_start_year"),
				r.PostFormValue("job_start_month"),
				r.PostFormValue("job_start_day"),
			),
			user.TZLocation())

		if err != nil {
			internalError(w, err)
			return
		}

		job.EmployerName = employerName
		job.StartTime = jobStart

		if job.ID > 0 {
			if err = job.Save(); err != nil {
				internalError(w, err)
				return
			}
		} else {
			daysPerYear, err := strconv.ParseFloat(r.PostFormValue("days_per_year"), 32)
			if err != nil {
				internalError(w, err)
				return
			}

			leaveYearStart, err := time.ParseInLocation("2006-1-2",
				fmt.Sprintf("%s-%s-%s",
					r.PostFormValue("leave_start_year"),
					r.PostFormValue("leave_start_month"),
					r.PostFormValue("leave_start_day"),
				),
				user.TZLocation())
			leaveYearEnd := leaveYearStart.AddDate(1, 0, -1)

			if err != nil {
				internalError(w, err)
				return
			}

			userNow := time.Now().In(user.TZLocation())
			if userNow.After(leaveYearEnd) || userNow.Before(leaveYearStart) {
				showError(w, http.StatusNotAcceptable, "Current leave year must not end before or start after today's date.")
				return
			}

			var charPool = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			exportSecret := make([]byte, 64)

			for i := range exportSecret {
				n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charPool))))
				if err != nil {
					internalError(w, err)
					return
				}
				exportSecret[i] = charPool[n.Int64()]
			}
			job.ExportSecret = string(exportSecret)

			job.LeaveAllowances = append(job.LeaveAllowances, model.LeaveAllowance{
				Minutes:      int32(daysPerYear * 24 * 60),
				Description:  "Annual leave allocation",
				StartTime:    leaveYearStart,
				EndTime:      leaveYearEnd,
				IsAdjustment: false,
			})

			user.Jobs = append(user.Jobs, job)
			if err = user.Save(); err != nil {
				internalError(w, err)
				return
			}

		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	output.HTML(w, http.StatusOK, "settings",
		&struct {
			FormValues url.Values
			JobID      uint64
			User       model.User
		}{
			formValues,
			job.ID,
			user,
		})
}
