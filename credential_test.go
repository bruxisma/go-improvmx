package improvmx

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"occult.work/doze/test"
)

type CredentialTestSuite struct {
	test.Suite
	session *Session
}

type CredentialErrorTestSuite struct {
	test.Suite
	session *Session
}

func (suite *CredentialTestSuite) SetupSuite() {
	router := test.NewRouter().
		Get(credentialsListPath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/credentials/")
			fmt.Fprintf(writer, `{
			"credentials": [
				{
					"created": 1581604970000,
					"usage": 0,
					"username": "richard"
				},
				{
					"created": 1581607028000,
					"usage": 0,
					"username": "monica"
				}
			],
			"success": true}`)
		}).
		Post(credentialsCreatePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/credentials/")
			fmt.Fprintf(writer, `{
				"credential": {
					"created": 1588236952000,
					"usage": 0,
					"username": "username"
				},
				"requires_new_mx_check": false,
				"success": true
			}`)
		}).
		Put(credentialsUpdatePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/credentials/username")
			fmt.Fprintf(writer, `{
				"credential": {
					"created": 1588236952000,
					"usage": 0,
					"username": "username"
				},
				"success": true
			}`)
		}).
		Delete(credentialsDeletePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/credentials/username")
			fmt.Fprintf(writer, `{ "success": true }`)
		})
	suite.Initialize(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func (suite *CredentialErrorTestSuite) SetupSuite() {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(420)
		fmt.Fprintf(writer, `{ "error": "fake error", "code": 420, "success": false }`)
	}
	router := test.NewRouter().
		Post(credentialsCreatePath, handler).
		Put(credentialsUpdatePath, handler).
		Get(credentialsListPath, handler)
	suite.Initialize(router)
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
	test.Run(t, new(CredentialTestSuite))
	test.Run(t, new(CredentialErrorTestSuite))
}

func (suite *CredentialTestSuite) TestList() {
	credentials, error := suite.session.Credentials.List(context.Background(), "example.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(credentials)
}

func (suite *CredentialTestSuite) TestCreate() {
	credential, error := suite.session.Credentials.Create(context.Background(), "example.com", User{
		Password: "password",
		Username: "username",
	})
	suite.Require().NoError(error)
	suite.Equal(credential.Username, "username")
}

func (suite *CredentialTestSuite) TestUpdate() {
	credential, error := suite.session.Credentials.Update(context.Background(), "example.com", User{
		Username: "username",
		Password: "password",
	})
	suite.Require().NoError(error)
	suite.Equal(credential.Username, "username")
}

func (suite *CredentialTestSuite) TestDelete() {
	error := suite.session.Credentials.Delete(context.Background(), "example.com", "username")
	suite.Require().NoError(error)
}

func (suite *CredentialErrorTestSuite) TestList() {
	credentials, error := suite.session.Credentials.List(context.Background(), "example.com")
	suite.Require().Error(error)
	suite.Require().Empty(credentials)
}

func (suite *CredentialErrorTestSuite) TestCreate() {
	credential, error := suite.session.Credentials.Create(context.Background(), "example.com", User{
		Password: "password",
		Username: "username",
	})
	suite.Require().Error(error)
	suite.Require().Empty(credential)
}

func (suite *CredentialErrorTestSuite) TestUpdate() {
	credential, error := suite.session.Credentials.Update(context.Background(), "example.com", User{
		Password: "password",
		Username: "username",
	})
	suite.Require().Error(error)
	suite.Require().Empty(credential)
}
