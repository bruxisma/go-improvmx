package improvmx

import (
	"embed"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

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

//func parseParameters(request *http.Request) (map[string]string, error) {
//	var parameters map[string]string
//	decoder := json.NewDecoder(request.Body)
//	error := decoder.Decode(&parameters)
//	return parameters, error
//}
//
//func badRequest(writer http.ResponseWriter, message string) {
//	writer.WriteHeader(http.StatusBadRequest)
//	fmt.Fprintf(writer, `{ "code": 400, "error": "400 Bad Request %s", "success": false }`, message)
//}

type CommonTestSuite struct {
	test.Suite
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

func (suite *CommonTestSuite) TestListOptionSetPage() {
	invalid := NewListOption().SetPage(-1)
	valid := NewListOption().SetPage(1)

	invaliderr := invalid.validate()
	validerr := valid.validate()

	suite.Require().Error(invaliderr)
	suite.Require().NoError(validerr)
}

func (suite *CommonTestSuite) TestListOptionSetLimit() {
	invalid := NewListOption().SetLimit(0)
	valid := NewListOption().SetLimit(50)

	invaliderr := invalid.validate()
	validerr := valid.validate()

	suite.Require().Error(invaliderr)
	suite.Require().NoError(validerr)
}
