package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (s *handlerTestSuite) TestGitHubOauthStateCookieSet() {
	req, err := http.NewRequest("GET", "https://localhost/", nil)
	if err != nil {
		s.T().Error(err)
	}

	w := httptest.NewRecorder()
	Index(w, req)

	assert.Regexp(s.T(), fmt.Sprintf("^%s=.*;", oauthStateCookieName), w.Header().Get("Set-Cookie"))
}
