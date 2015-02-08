package handler

import "net/http"

func Dashboard(w http.ResponseWriter, r *http.Request) {
	output.HTML(w, http.StatusOK, "dashboard", nil)
}
