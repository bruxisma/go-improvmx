package improvmx

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/suite"
	"occult.work/improvmx/internal"
)

type DomainTestSuite struct {
	internal.Suite
	session *Session
}

func (suite *DomainTestSuite) SetupSuite() {
	router := internal.NewRouter().
		Get(domainVerifyPath, suite.FileResponseHandler("testdata/domain/verify.json")).
		Get(domainListPath, suite.FileResponseHandler("testdata/domain/list.json")).
		Get(domainLogsPath, suite.FileResponseHandler("testdata/domain/logs.json")).
		Get(domainReadPath, suite.FileResponseHandler("testdata/domain/read.json"))
	suite.InitializeWithRouter(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func TestDomain(t *testing.T) {
	suite.Run(t, new(DomainTestSuite))
}

func (suite *DomainTestSuite) TestList() {
	domains, error := suite.session.Domains.List(context.Background())
	suite.Require().NoError(error)
	suite.Require().NotEmpty(domains)
}

func (suite *DomainTestSuite) TestLogs() {
	logs, error := suite.session.Domains.Logs(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(logs)
}

func (suite *DomainTestSuite) TestCreate() {
	suite.Skip()
}

func (suite *DomainTestSuite) TestRead() {
	domain, error := suite.session.Domains.Read(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
	suite.False(domain.Active)
	suite.NotEmpty(domain.Aliases)
}

func (suite *DomainTestSuite) TestUpdate() {
	suite.Skip()
}

func (suite *DomainTestSuite) TestDelete() {
	suite.Skip()
}

func (suite *DomainTestSuite) TestVerify() {
	error := suite.session.Domains.Verify(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
}
