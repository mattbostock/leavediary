package handler

import "net/http"

func showError(w http.ResponseWriter, msg string, statusCode int) {
	output.HTML(w, statusCode, "error", msg)
}
