package improvmx

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"occult.work/improvmx/internal"
)

type CredentialTestSuite struct {
	internal.Suite
	session *Session
}

func (suite *CredentialTestSuite) SetupSuite() {
	router := internal.NewRouter().
		Delete(credentialsDeletePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/credentials/username")
			fmt.Fprintf(writer, `{ "success": true }`)
		})
	suite.InitializeWithRouter(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func TestCredentialPaths(t *testing.T) {
	assert := assert.New(t)
	fullpath := regexp.MustCompile("/$")
	resource := regexp.MustCompile("[^/]$")

	assert.Regexp(fullpath, credentialsListPath)
	assert.Regexp(fullpath, credentialsCreatePath)
	assert.Regexp(resource, credentialsUpdatePath)
	assert.Regexp(resource, credentialsDeletePath)
}

func TestCredential(t *testing.T) {
	suite.Run(t, new(CredentialTestSuite))
}

func (suite *CredentialTestSuite) TestList() {
	suite.Skip()
}

func (suite *CredentialTestSuite) TestCreate() {
	suite.Skip()
}

func (suite *CredentialTestSuite) TestUpdate() {
	suite.Skip()
}

func (suite *CredentialTestSuite) TestDelete() {
	error := suite.session.Credentials.Delete(context.Background(), "example.com", "username")
	suite.Require().NoError(error)
}
