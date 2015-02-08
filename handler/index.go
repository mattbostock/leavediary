package handler

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	output.HTML(w, http.StatusOK, "index", nil)
}
