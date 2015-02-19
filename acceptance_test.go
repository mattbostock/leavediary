package main

import (
	"fmt"
	"testing"

	agouti "github.com/sclevine/agouti/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var baseURL string

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

	baseURL = fmt.Sprintf("http://%s", config.addr)

	config.gitHubClientID = "abc"
	config.gitHubClientSecret = "xyz"
	go main()

	s.driver, err = agouti.PhantomJS()
	s.driver.Start()
	s.page, err = s.driver.Page(agouti.Use().Browser("chrome"))

	if err != nil {
		s.T().Error(err)
	}
}

func (s *acceptanceTestSuite) TearDownSuite() {
	s.driver.Stop()
}

func (s *acceptanceTestSuite) TestDebugVarsExposed() {
	_ = s.page.Navigate(baseURL + "/debug/vars")
	bodyText, _ := s.page.Find("body").Text()

	assert.Contains(s.T(), bodyText, "cmdline")
	assert.Contains(s.T(), bodyText, "memstats")
}

func (s *acceptanceTestSuite) TestHomePageForJavascriptErrors() {
	_ = s.page.Navigate(baseURL)
	logs, _ := s.page.ReadLogs("browser", true)

	for _, log := range logs {
		assert.NotEqual(s.T(), "WARNING", log.Level, log.Message)
		assert.NotEqual(s.T(), "SEVERE", log.Level, log.Message)
	}
}
