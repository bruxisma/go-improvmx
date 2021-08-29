package improvmx

import (
	"context"
	_ "embed"
	"testing"

	"occult.work/doze/test"
)

type AccountTestSuite struct {
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

func TestAccount(t *testing.T) {
	test.Run(t, new(AccountTestSuite))
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
