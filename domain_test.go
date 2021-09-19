package improvmx

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"occult.work/doze/test"
)

type DomainTestSuite struct {
	test.Suite
	session *Session
}

func (suite *DomainTestSuite) SetupSuite() {
	router := test.NewRouter().
		Get(domainVerifyPath, suite.FileResponseHandler("testdata/domain/verify.json")).
		Get(domainListPath, suite.FileResponseHandler("testdata/domain/list.json")).
		Get(domainLogsPath, suite.FileResponseHandler("testdata/domain/logs.json")).
		Get(domainReadPath, suite.FileResponseHandler("testdata/domain/read.json")).
		Post(domainCreatePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/")
			domain := Domain{
				Active: true, Name: "example.com"}
			data, error := json.Marshal(domain)
			suite.Require().NoError(error)
			fmt.Fprintf(writer, `{ "domain": %s, "success": true }`, string(data))
		}).
		Put(domainUpdatePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/piedpiper.com/")
			option := &DomainOption{}
			bytes, error := ioutil.ReadAll(request.Body)
			suite.Require().NoError(error)
			error = json.Unmarshal(bytes, option)
			suite.Require().NoError(error)

			data, error := suite.Data.ReadFile("testdata/domain/read.json")
			suite.Require().NoError(error)
			suite.Require().NotEmpty(data)

			response := domainResponse{}
			error = json.Unmarshal(data, &response)
			suite.Require().NoError(error)
			domain := response.Domain

			if option.Label != "" {
				domain.Whitelabel = option.Label
			}
			if option.Email != "" {
				domain.NotificationEmail = option.Email
			}

			data, error = json.Marshal(domain)

			suite.Require().NoError(error)
			suite.Require().NotEmpty(data)

			fmt.Fprintf(writer, `{ "domain": %s, "success": true }`, string(data))
		}).
		Delete(domainDeletePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/")
			fmt.Fprintf(writer, `{ "success": true }`)
		})
	suite.Initialize(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func TestDomain(t *testing.T) {
	test.Run(t, new(DomainTestSuite))
}

func (suite *DomainTestSuite) TestList() {
	domains, error := suite.session.Domains.List(context.Background())
	suite.Require().NoError(error)
	suite.Require().NotEmpty(domains)
}

func (suite *DomainTestSuite) TestListWithOptions() {
	options := NewListOption().
		SetIsActive(true).
		SetStartsWith("n/a").
		SetLimit(20).
		SetPage(1)
	domains, error := suite.session.Domains.List(context.Background(), options)
	suite.Require().NoError(error)
	suite.Require().NotEmpty(domains)
}

func (suite *DomainTestSuite) TestListWithInvalidOptions() {
	options := NewListOption().SetPage(-1)
	domains, error := suite.session.Domains.List(context.Background(), options)
	suite.Require().Error(error)
	suite.Require().Empty(domains)
}

func (suite *DomainTestSuite) TestLogs() {
	logs, error := suite.session.Domains.Logs(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(logs)
}

func (suite *DomainTestSuite) TestCreate() {
	domain, error := suite.session.Domains.Create(context.Background(), "example.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(domain)
}

func (suite *DomainTestSuite) TestRead() {
	domain, error := suite.session.Domains.Read(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(domain)
	suite.False(domain.Active)
	suite.NotEmpty(domain.Aliases)
}

func (suite *DomainTestSuite) TestUpdate() {
	domain, error := suite.session.Domains.Update(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
	suite.NotEmpty(domain)
}

func (suite *DomainTestSuite) TestDelete() {
	error := suite.session.Domains.Delete(context.Background(), "example.com")
	suite.Require().NoError(error)
}

func (suite *DomainTestSuite) TestVerify() {
	error := suite.session.Domains.Verify(context.Background(), "piedpiper.com")
	suite.Require().NoError(error)
}
