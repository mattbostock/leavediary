package handler

import (
	"net/http"
	"net/http/httptest"

	"github.com/mattbostock/leavediary/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func (s *handlerTestSuite) TestOauthCallbackStateMismatch() {
	r, err := http.NewRequest("GET", "https://localhost/?state=1234", nil)
	if err != nil {
		s.T().Error(err)
	}

	r.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "4321"})
	w := httptest.NewRecorder()

	GithubOauthCallback(w, r)

	assert.Equal(s.T(), "GitHub state mismatch during Oauth callback", s.log.LastEntry().Message)
	assert.Equal(s.T(), http.StatusTemporaryRedirect, w.Code)
	assert.Equal(s.T(), "/", w.Header().Get("Location"))
}

func (s *handlerTestSuite) TestOauthExchangeFailed() {
	r, err := http.NewRequest("GET", "https://localhost/?state=1234", nil)
	if err != nil {
		s.T().Error(err)
	}

	r.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "1234"})
	w := httptest.NewRecorder()

	// Use an Oauth config that should fail the Oauth exchange
	oauthConfig.Endpoint = oauth2.Endpoint{"https://127.0.0.1/404", "https://127.0.0.1/404"}

	GithubOauthCallback(w, r)

	assert.Contains(s.T(), s.log.LastEntry().Message, "GitHub Oauth exchange failed")
	assert.Equal(s.T(), http.StatusTemporaryRedirect, w.Code)
	assert.Equal(s.T(), "/", w.Header().Get("Location"))
}

func (s *handlerTestSuite) TestGitHubAuthenticationNewUser() {
	r, err := http.NewRequest("GET", "https://localhost/?state=1234", nil)
	if err != nil {
		s.T().Error(err)
	}

	r.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "1234"})
	w := httptest.NewRecorder()

	setupTestDB()

	GithubOauthCallback(w, r)

	u, err := model.FindUser(1)
	if err != nil {
		panic(err)
	}

	assert.Equal(s.T(), "octocat@github.com", u.Email)
	assert.Equal(s.T(), "monalisa octocat", u.Name)
	assert.Equal(s.T(), 1, u.GitHubID)

	assert.Equal(s.T(), http.StatusTemporaryRedirect, w.Code)
	assert.Equal(s.T(), "/dashboard", w.Header().Get("Location"))

	assert.Len(s.T(), s.log.Entries, 0)
}
