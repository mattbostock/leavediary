package handler

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/mattbostock/timeoff/handler/mocks/logrus"
	"github.com/mattbostock/timeoff/middleware/sessions"
	"github.com/mattbostock/timeoff/model"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type handlerTestSuite struct {
	suite.Suite
	log *mockhook.Mockhook
}

func TestHandlerTests(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

func (s *handlerTestSuite) SetupTest() {
	log = logrus.New()
	s.log = &mockhook.Mockhook{}
	log.Hooks.Add(s.log)

	oauthConfig = &oauth2.Config{
		ClientID:     "abc",
		ClientSecret: "xyz",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
	sessionManager = sessions.New("test_session", securecookie.GenerateRandomKey(32))
}

func setupTestDB() {
	dbDialect := os.Getenv("DB_DIALECT")
	dbDataSource := os.Getenv("DB_DATASOURCE")

	if dbDialect == "" && dbDataSource == "" {
		dbDialect = "sqlite3"
		dbDataSource = ":memory:"
	}

	model.SetLogger(log)
	model.InitDB(dbDialect, dbDataSource)
}
