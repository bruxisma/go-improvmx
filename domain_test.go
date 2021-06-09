package improvmx

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
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
		Get(domainReadPath, suite.FileResponseHandler("testdata/domain/read.json")).
		Delete(domainDeletePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/")
			fmt.Fprintf(writer, `{ "success": true }`)
		})
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
	error := suite.session.Domains.Delete(context.Background(), "example.com")
	suite.Require().NoError(error)
}

func (suite *DomainTestSuite) TestVerify() {
	error := suite.session.Domains.Verify(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
}
