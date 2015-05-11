package handler

import (
	"encoding/csv"
	"net/http"
)

func ExportCSV(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	if user.ID == 0 { // no current user; not logged in
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	job := user.Jobs[0] // only support one job for now
	requests, err := job.RequestsLastYearAndFuture()
	if err != nil {
		internalError(w, err)
		return
	}

	csv := csv.NewWriter(w)

	record := make([]string, 4)

	// write CSV headers
	record = []string{"Name", "Start time", "End time", "Description"}
	if err := csv.Write(record); err != nil {
		internalError(w, err)
		return
	}

	for _, req := range requests {
		record = []string{user.Name, req.StartTime.String(), req.EndTime.String(), req.Description}
		if err := csv.Write(record); err != nil {
			internalError(w, err)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\"timeoff.csv\"")
	w.Header().Set("Content-Type", "text/csv")
	csv.Flush()
}
