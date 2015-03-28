package handler

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func (s *handlerTestSuite) TestTooManyRequests() {
	r, err := http.NewRequest("GET", "https://localhost/", nil)
	if err != nil {
		s.T().Error(err)
	}

	w := httptest.NewRecorder()
	w.Header().Set("Retry-After", "11")

	r.RemoteAddr = "127.0.0.2:55555"

	TooManyRequests(w, r)

	assert.Equal(s.T(), 429, w.Code)
	assert.Equal(s.T(), "Too many requests from your IP 127.0.0.2; try again in 11 seconds\n", w.Body.String())
}
