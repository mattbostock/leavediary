package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mattbostock/leavediary/model"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(user.Jobs) == 0 {
		log.Infoln("No jobs defined")
		http.Redirect(w, r, "/settings", http.StatusTemporaryRedirect)
		return
	}

	job := user.CurrentJob()
	currentLeaveAllowance, err := job.CurrentLeaveAllowance()
	if err != nil {
		internalError(w, err)
		return
	}
	if currentLeaveAllowance.ID == 0 {
		log.Infoln("No current leave allowance")
		http.Redirect(w, r, "/allowance", http.StatusTemporaryRedirect)
		return
	}

	remainingMinutes, err := currentLeaveAllowance.RemainingTime()
	if err != nil {
		internalError(w, err)
		return
	}
	remainingDays := float32(remainingMinutes) / 60 / 24

	var nextOnLeave model.LeaveRequest
	var pastRequests, upcomingRequests []model.LeaveRequest

	requests, err := job.RequestsLastYearAndFuture()
	if err != nil {
		internalError(w, err)
		return
	}

	for _, req := range requests {
		if req.Before(time.Now().In(user.TZLocation())) {
			pastRequests = append(pastRequests, req)
		}
	}
	for _, req := range requests {
		if req.After(time.Now().AddDate(0, 0, -1).In(user.TZLocation())) {
			upcomingRequests = append(upcomingRequests, req)
		}
	}

	if len(upcomingRequests) > 0 {
		nextOnLeave = upcomingRequests[0]
	}

	calendarURL := &url.URL{
		Scheme: "https",
		Host:   r.Host,
		Path:   fmt.Sprintf("/export/ics/%s", job.ExportSecret),
	}

	webCalURLNoScheme := &url.URL{}
	*webCalURLNoScheme = *calendarURL
	webCalURLNoScheme.Scheme = ""

	output.HTML(w, http.StatusOK, "dashboard", &struct {
		CurrentLeaveAllowance model.LeaveAllowance
		CalendarURL           *url.URL
		NextOnLeave           model.LeaveRequest
		PastRequests          []model.LeaveRequest
		UpcomingRequests      []model.LeaveRequest
		RemainingDays         float32
		User                  model.User
		WebCalURLNoScheme     *url.URL
	}{
		currentLeaveAllowance,
		calendarURL,
		nextOnLeave,
		pastRequests,
		upcomingRequests,
		remainingDays,
		user,
		webCalURLNoScheme,
	})
}
