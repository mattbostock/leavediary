package handler

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	o.HTML(w, http.StatusOK, "index", "world")
}
