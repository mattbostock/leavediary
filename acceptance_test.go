package main

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"testing"

	agouti "github.com/sclevine/agouti/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type acceptanceTestSuite struct {
	suite.Suite
	driver agouti.WebDriver
	page   agouti.Page
}

func TestAcceptanceTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance tests in short mode.")
	}
	suite.Run(t, new(acceptanceTestSuite))
}

func (s *acceptanceTestSuite) SetupSuite() {
	var err error

	s.driver, err = agouti.PhantomJS()
	s.driver.Start()
	s.page, err = s.driver.Page(agouti.Use().Browser("chrome"))

	if err != nil {
		s.T().Error(err)
	}

	go main()
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.driver.Stop()
}

func (s *acceptanceTestSuite) TestDebugVarsExposed() {
	_ = s.page.Navigate("http://localhost:3000/debug/vars")
	bodyText, _ := s.page.Find("body").Text()

	assert.Contains(s.T(), bodyText, "cmdline")
	assert.Contains(s.T(), bodyText, "memstats")
}

func (s *acceptanceTestSuite) TestHomePageForJavascriptErrors() {
	_ = s.page.Navigate("http://localhost:3000/")
	logs, _ := s.page.ReadLogs("browser", true)

	for _, log := range logs {
		assert.NotEqual(s.T(), log.Level, "WARNING", log.Message)
		assert.NotEqual(s.T(), log.Level, "SEVERE", log.Message)
	}
}

func (s *acceptanceTestSuite) TestGZIPEnabledWhenSupported() {
	for _, encoding := range []string{"", "gzip"} {
		req, _ := http.NewRequest("GET", "http://localhost:3000", nil)
		req.Header.Add("Accept-Encoding", encoding)

		resp, _ := http.DefaultTransport.RoundTrip(req)
		defer resp.Body.Close()

		var body []byte

		if encoding == "gzip" {
			gzBody, _ := gzip.NewReader(resp.Body)
			defer gzBody.Close()
			body, _ = ioutil.ReadAll(gzBody)
		} else {
			body, _ = ioutil.ReadAll(resp.Body)
		}

		assert.Equal(s.T(), resp.Header.Get("Content-Encoding"), encoding)
		assert.Contains(s.T(), string(body), "TimeOff")
	}
}
