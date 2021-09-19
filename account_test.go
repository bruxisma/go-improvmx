package improvmx

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"occult.work/doze/test"
)

type AccountTestSuite struct {
	test.Suite
	session *Session
}

type AccountErrorTestSuite struct {
	test.Suite
	session *Session
}

func (suite *AccountTestSuite) SetupSuite() {
	router := test.NewRouter().
		Get(accountLabelsPath, suite.FileResponseHandler("testdata/account/labels.json")).
		Get(accountReadPath, suite.FileResponseHandler("testdata/account/read.json"))

	suite.Initialize(router)
	suite.Data = &testData
	suite.session, _ = New("account-test-suite", WithBaseURL(suite.Server.URL))
}

func (suite *AccountErrorTestSuite) SetupSuite() {
	router := test.NewRouter().
		Get(accountLabelsPath, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(420)
			fmt.Fprintf(writer, `{ "error": "fake error", "code": 420, "success": false }`)
		}).
		Get(accountReadPath, func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(420)
			fmt.Fprintf(writer, `{ "error": "fake error", "code": 420, "success": false }`)
		})
	suite.Initialize(router)
	suite.Data = &testData
	suite.session, _ = New("account-error-test-suite", WithBaseURL(suite.Server.URL))
}

func TestAccount(t *testing.T) {
	test.Run(t, new(AccountTestSuite))
	test.Run(t, new(AccountErrorTestSuite))
}

func (suite *AccountTestSuite) TestRead() {
	account, error := suite.session.Account.Read(context.Background())
	suite.Require().NoError(error)
	suite.Equal("1234", account.Last4)
	suite.Equal("US", account.Country)
	suite.Equal(10, account.Limits.RateLimit)
}

func (suite *AccountTestSuite) TestLabels() {
	labels, error := suite.session.Account.Labels(context.Background())
	suite.Require().NoError(error)
	suite.Require().NotEmpty(labels)
	suite.Equal(labels[0].Name, "example.com")
}

func (suite *AccountErrorTestSuite) TestRead() {
	account, error := suite.session.Account.Read(context.Background())
	suite.Require().Error(error)
	suite.Require().Empty(account)
}

func (suite *AccountErrorTestSuite) TestLabels() {
	labels, error := suite.session.Account.Labels(context.Background())
	suite.Require().Error(error)
	suite.Require().Empty(labels)
}
