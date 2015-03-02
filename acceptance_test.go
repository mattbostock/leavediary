package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
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

	baseURL = fmt.Sprintf("https://%s/", config.addr)

	config.gitHubClientID = "abc"
	config.gitHubClientSecret = "xyz"
	config.tlsCert = filepath.Join(fixturesPath, "cert.pem")
	config.tlsKey = filepath.Join(fixturesPath, "key.pem")
	// FIXME need to configure and clean up database

	go main()

	s.driver = agouti.NewWebDriver("http://{{.Address}}", []string{"phantomjs", "--webdriver={{.Address}}", "--ignore-ssl-errors=true"})
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
	err = waitForServer(baseURL)
	if err != nil {
		s.T().Error(err)
	}
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.driver.Stop()
}

func (s *acceptanceTestSuite) TestDebugVarsExposed() {
	testURL := baseURL + "debug/vars"
	resp, err := http.Get(testURL)
	if err != nil {
		s.T().Error(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		s.T().Error(err)
	}

	assert.Contains(s.T(), string(body), "cmdline")
	assert.Contains(s.T(), string(body), "memstats")
}

func (s *acceptanceTestSuite) TestHomePageForJavascriptErrors() {
	testURL := baseURL
	err := s.page.Navigate(testURL)
	if err != nil {
		s.T().Error(err)
	}

	logs, err := s.page.ReadAllLogs("browser")
	if err != nil {
		s.T().Error(err)
	}

	u, _ := s.page.URL()
	assert.Equal(s.T(), testURL, u)

	for _, log := range logs {
		assert.NotEqual(s.T(), "WARNING", log.Level, log.Message)
		assert.NotEqual(s.T(), "SEVERE", log.Level, log.Message)
	}
}

func (s *acceptanceTestSuite) TestPageNotFound() {
	resp, err := http.Get(baseURL + "non-existentent-page")
	if err != nil {
		s.T().Error(err)
	}

	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

func waitForServer(u string) error {
	for {
		_, err := http.Get(u)

		if err == nil {
			break
		} else {
			switch err.(type) {
			case *url.Error:
				if strings.HasSuffix(err.Error(), "connection refused") {
					time.Sleep(500 * time.Millisecond)
					continue
				}
				return err
			default:
				return err
			}
		}
	}

	return nil
}
