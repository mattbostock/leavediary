package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/mattbostock/leavediary/handler/mocks/logrus"
	"github.com/mattbostock/leavediary/middleware/sessions"
	"github.com/mattbostock/leavediary/model"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
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
	oauthServer := mockGithubServer()

	var err error
	githubAPIBaseURL, err = url.Parse(oauthServer.URL)
	if err != nil {
		panic(err)
	}

	oauthConfig = &oauth2.Config{
		ClientID:     "abc",
		ClientSecret: "xyz",
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauthServer.URL + "/foo",
			TokenURL: oauthServer.URL + "/login/oauth/access_token",
		},
		Scopes: []string{"user:email"},
	}
	sessionManager = sessions.New("test_session", securecookie.GenerateRandomKey(32))
}

func mockGithubServer() *httptest.Server {
	// FIXME Split this into an Oauth server and a GitHub API server at some point
	m := http.NewServeMux()
	server := httptest.NewServer(m)

	var accessTokenPayload = []byte(`access_token=sekret&scope=user%3Aemail&token_type=bearer`)
	var userPayload = []byte(`
{
  "login": "octocat",
  "id": 1,
  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
  "gravatar_id": "",
  "url": "https://api.github.com/users/octocat",
  "html_url": "https://github.com/octocat",
  "followers_url": "https://api.github.com/users/octocat/followers",
  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
  "organizations_url": "https://api.github.com/users/octocat/orgs",
  "repos_url": "https://api.github.com/users/octocat/repos",
  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
  "received_events_url": "https://api.github.com/users/octocat/received_events",
  "type": "User",
  "site_admin": false,
  "name": "monalisa octocat",
  "company": "GitHub",
  "blog": "https://github.com/blog",
  "location": "San Francisco",
  "email": "octocat@github.com",
  "hireable": false,
  "bio": "There once was...",
  "public_repos": 2,
  "public_gists": 1,
  "followers": 20,
  "following": 0,
  "created_at": "2008-01-14T04:33:35Z",
  "updated_at": "2008-01-14T04:33:35Z",
  "total_private_repos": 100,
  "owned_private_repos": 100,
  "private_gists": 81,
  "disk_usage": 10000,
  "collaborators": 8,
  "plan": {
    "name": "Medium",
    "space": 400,
    "private_repos": 20,
    "collaborators": 0
  }
}
`)
	var userEmailsPayload = []byte(`
[
  {
    "email": "octocat@github.com",
    "verified": true,
    "primary": true
  }
]
`)

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login/oauth/access_token":
			w.Write(accessTokenPayload)
			return
		case "/user":
			w.Write(userPayload)
			return
		case "/user/emails":
			w.Write(userEmailsPayload)
			return
		default:
			panic("Request received at unexpected endpoint: " + r.URL.Path)
		}
	})

	return server
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
