package handler

import "net/http"

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionManager.Logout(w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
