package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sclevine/agouti"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const fixturesPath = "test_fixtures"

var baseURL string

type acceptanceTestSuite struct {
	suite.Suite
	driver *agouti.WebDriver
	page   *agouti.Page
}

func TestAcceptanceTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance tests in short mode.")
	}
	suite.Run(t, new(acceptanceTestSuite))
}

func (s *acceptanceTestSuite) SetupSuite() {
	var err error

	baseURL = fmt.Sprintf("https://%s", config.addr)

	config.gitHubClientID = "abc"
	config.gitHubClientSecret = "xyz"
	config.tlsCert = filepath.Join(fixturesPath, "cert.pem")
	config.tlsKey = filepath.Join(fixturesPath, "key.pem")

	go main()

	s.driver = agouti.PhantomJS()
	s.driver.Start()

	s.page, err = s.driver.NewPage(agouti.Desired(agouti.NewCapabilities().Browser("chrome")))
	if err != nil {
		s.T().Error(err)
	}

	// don't verify our development TLS certificates
	http.DefaultTransport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Make sure server is actually running (see go main() call above) before
	// running tests to avoid race conditions
	for {
		_, err := http.Get(baseURL)

		if err == nil {
			break
		} else {
			switch err.(type) {
			case *url.Error:
				if strings.HasSuffix(err.Error(), "connection refused") {
					time.Sleep(500 * time.Millisecond)
					continue
				}
			default:
				s.T().Error(err)
				break
			}
		}
	}
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.driver.Stop()
}

func (s *acceptanceTestSuite) TestDebugVarsExposed() {
	err := s.page.Navigate(baseURL + "/debug/vars")
	if err != nil {
		s.T().Error(err)
	}

	bodyText, err := s.page.Find("body").Text()
	if err != nil {
		s.T().Error(err)
	}

	assert.Contains(s.T(), bodyText, "cmdline")
	assert.Contains(s.T(), bodyText, "memstats")
}

func (s *acceptanceTestSuite) TestHomePageForJavascriptErrors() {
	err := s.page.Navigate(baseURL)
	if err != nil {
		s.T().Error(err)
	}

	logs, err := s.page.ReadAllLogs("browser")
	if err != nil {
		s.T().Error(err)
	}

	for _, log := range logs {
		assert.NotEqual(s.T(), "WARNING", log.Level, log.Message)
		assert.NotEqual(s.T(), "SEVERE", log.Level, log.Message)
	}
}

func (s *acceptanceTestSuite) TestPageNotFound() {
	resp, err := http.Get(baseURL + "/non-existentent-page")
	if err != nil {
		s.T().Error(err)
	}

	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}
