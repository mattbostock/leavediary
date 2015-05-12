package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mattbostock/leavediary/model"
	"github.com/soh335/ical"
)

func ExportICS(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get(":secret")
	job, err := model.FindJobFromExportSecret(secret)
	if err != nil {
		internalError(w, err)
		return
	}

	if job.ID == 0 {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user, err := model.FindUser(job.UserID)
	if err != nil {
		internalError(w, err)
		return
	}

	requests, err := job.RequestsLastYearAndFuture()
	if err != nil {
		internalError(w, err)
		return
	}

	cal := ical.NewBasicVCalendar()
	// PRODID must comply with RFC5545
	cal.PRODID = fmt.Sprintf("Lucky Llama Ltd: LeaveDiary v%s", version)
	cal.X_WR_TIMEZONE = user.TZLocation().String()
	cal.X_WR_CALNAME = fmt.Sprintf("Time off for %s", user.Name)

	for _, e := range requests {
		v := &ical.VEvent{
			UID:     "LeaveDiary-" + strconv.FormatUint(e.ID, 10),
			DTSTAMP: e.UpdatedAt,
			DTSTART: e.StartTime,
			DTEND:   e.EndTime,
			SUMMARY: user.Name + " annual leave: " + e.Description,
			TZID:    e.StartTime.Location().String(),
		}
		cal.VComponent = append(cal.VComponent, v)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\"leavediary.ics\"")
	w.Header().Set("Content-Type", "text/calendar")
	cal.Encode(w)
}
