package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/mattbostock/leavediary/model"
)

func Index(w http.ResponseWriter, r *http.Request) {
	randStr := make([]byte, 12)
	_, err := rand.Read(randStr)
	if err != nil {
		internalError(w, err)
		return
	}
	state := base64.StdEncoding.EncodeToString(randStr)

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   24 * 60 * 60,
		Secure:   true,
		HttpOnly: true,
	})

	url := oauthConfig.AuthCodeURL(state)
	user := currentUser(r)

	output.HTML(w, http.StatusOK, "index", &struct {
		GitHubOauthURL string
		User           model.User
	}{url, user})
}
