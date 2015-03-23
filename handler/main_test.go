package handler

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/mattbostock/timeoff/middleware/sessions"
	"github.com/stretchr/testify/suite"
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
