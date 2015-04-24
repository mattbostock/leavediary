package handler

import (
	"net/http"

	"github.com/google/go-github/github"
	"github.com/mattbostock/timeoff/model"
	"golang.org/x/oauth2"
)

func GithubOauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	stateCookie, _ := r.Cookie(oauthStateCookieName)

	if state != stateCookie.Value {
		log.Errorln("GitHub state mismatch during Oauth callback")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   oauthStateCookieName,
		MaxAge: -1,
	})

	t, err := oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Errorln("GitHub Oauth exchange failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	httpClient := oauthConfig.Client(oauth2.NoContext, t)
	githubClient := github.NewClient(httpClient)

	// set this explicitly as we'll override it in tests
	githubClient.BaseURL = githubAPIBaseURL

	user, _, err := githubClient.Users.Get("")

	if err != nil {
		log.Infoln("Authentication failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	u := &model.User{
		Name:     *user.Name,
		GitHubID: uint64(*user.ID),
	}

	emails, _, err := githubClient.Users.ListEmails(nil)
	if err != nil {
		log.Errorln(err)
	}

	for _, e := range emails {
		if *e.Primary && *e.Verified {
			u.Email = *e.Email
			break
		}
	}

	u.UpdateOrCreate()

	err = sessionManager.SetCookie(w, uint64(u.ID))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}
