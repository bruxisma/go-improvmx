package improvmx

import (
	"embed"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"occult.work/doze/test"
)

//go:embed testdata
var testData embed.FS

func setupSession(server *httptest.Server) *Session {
	session, _ := New(os.Getenv("IMPROVMX_API_TOKEN"), WithBaseURL(server.URL))
	if session.client.HostURL != server.URL {
		return nil
	}
	return session
}

type CommonTestSuite struct {
	test.Suite
}

func TestError(t *testing.T) {
	assert := assert.New(t)
	error := Error{Message: "text", Code: 420}
	assert.Equal(error.Error(), `420: text`)
}

func TestCommon(t *testing.T) {
	test.Run(t, new(CommonTestSuite))
}

// Ensure that this function matches the const path values found in common.go
func (suite *CommonTestSuite) TestPaths() {
	regex := regexp.MustCompile("/$")
	suite.Regexp(regex, accountLabelsPath)
	suite.Regexp(regex, accountReadPath)

	suite.Regexp(regex, aliasListPath)
	suite.Regexp(regex, aliasCreatePath)
	suite.Regexp(regex, aliasReadPath)
	suite.Regexp(regex, aliasUpdatePath)
	suite.Regexp(regex, aliasDeletePath)
	suite.Regexp(regex, aliasLogsPath)

	suite.Regexp(regex, domainListPath)
	suite.Regexp(regex, domainCreatePath)
	suite.Regexp(regex, domainReadPath)
	suite.Regexp(regex, domainUpdatePath)
	suite.Regexp(regex, domainDeletePath)
	suite.Regexp(regex, domainVerifyPath)
	suite.Regexp(regex, domainLogsPath)
}

func (suite *CommonTestSuite) TestListOptionSetStartsWith() {
	option := NewListOption().SetStartsWith("test")
	suite.Require().Equal(option.startsWith, "test")
}

func (suite *CommonTestSuite) TestListOptionSetIsActive() {
	yes := NewListOption().SetIsActive(true)
	no := NewListOption().SetIsActive(false)
	suite.Require().NotEmpty(yes)
	suite.Require().NotEmpty(no)
	suite.Require().Equal(*yes.isActive, 1)
	suite.Require().Equal(*no.isActive, 0)
}

func (suite *CommonTestSuite) TestListOptionSetLimit() {
	invalid := NewListOption().SetLimit(0)
	valid := NewListOption().SetLimit(50)

	invaliderr := invalid.validate()
	validerr := valid.validate()

	suite.Require().Error(invaliderr)
	suite.Require().NoError(validerr)
}

func (suite *CommonTestSuite) TestListOptionSetPage() {
	invalid := NewListOption().SetPage(-1)
	valid := NewListOption().SetPage(1)

	invaliderr := invalid.validate()
	validerr := valid.validate()

	suite.Require().Error(invaliderr)
	suite.Require().NoError(validerr)
}

func (suite *CommonTestSuite) TestListOptionClearStartsWith() {
	option := NewListOption().SetStartsWith("test").ClearStartsWith()
	suite.Require().Equal(option.startsWith, "")
}

func (suite *CommonTestSuite) TestListOptionClearIsActive() {
	option := NewListOption().SetIsActive(true).ClearIsActive()
	suite.Require().Empty(option.isActive)
}

func (suite *CommonTestSuite) TestListOptionClearLimit() {
	option := NewListOption().SetLimit(50).ClearLimit()
	suite.Require().Empty(option.limit)
}

func (suite *CommonTestSuite) TestListOptionClearPage() {
	option := NewListOption().SetPage(-1).ClearPage()
	suite.Require().Empty(option.page)
}
