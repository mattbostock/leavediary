package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/mattbostock/timeoff/middleware/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type handlerTestSuite struct {
	suite.Suite
}

func TestHandlerTests(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (s *handlerTestSuite) SetupSuite() {
	log = logrus.New()
	oauthConfig = &oauth2.Config{
		ClientID:     "abc",
		ClientSecret: "xyz",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
	sessionManager = sessions.New("test_session", securecookie.GenerateRandomKey(32))
}

func (s *handlerTestSuite) TestGitHubOauthStateCookieSet() {
	req, err := http.NewRequest("GET", "https://localhost/", nil)
	if err != nil {
		s.T().Error(err)
	}

	w := httptest.NewRecorder()
	Index(w, req)

	assert.Regexp(s.T(), fmt.Sprintf("^%s=.*;", oauthStateCookieName), w.Header().Get("Set-Cookie"))
}
