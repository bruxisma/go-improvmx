package improvmx

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"occult.work/doze/test"
)

type DomainTestSuite struct {
	test.Suite
	session *Session
}

type DomainErrorTestSuite (DomainTestSuite)

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

func (suite *DomainErrorTestSuite) SetupSuite() {
	handler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(420)
		fmt.Fprint(writer, `{ "error": "fake error", code: 420, "success": false }`)
	}
	router := test.NewRouter().
		Get(domainListPath, handler).
		Get(domainLogsPath, handler).
		Get(domainReadPath, handler).
		Post(domainCreatePath, handler).
		Put(domainUpdatePath, handler)
	suite.Initialize(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func TestDomain(t *testing.T) {
	test.Run(t, new(DomainTestSuite))
	test.Run(t, new(DomainErrorTestSuite))
}

func TestDomainOption(test *testing.T) {
	assert := assert.New(test)
	option := getDomainOption(
		DomainOption{Label: "1st", Email: "1@example.com"},
		DomainOption{Label: "2nd", Email: "2@example.com"},
	)
	assert.Equal(option.Label, "1st")
	assert.Equal(option.Email, "1@example.com")
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

func (suite *DomainErrorTestSuite) TestLogs() {
	logs, error := suite.session.Domains.Logs(context.Background(), "example.com")
	suite.Require().Error(error)
	suite.Require().Empty(logs)
}

func (suite *DomainErrorTestSuite) TestList() {
	domains, error := suite.session.Domains.List(context.Background())
	suite.Require().Error(error)
	suite.Require().Empty(domains)
}

func (suite *DomainErrorTestSuite) TestCreate() {
	domain, error := suite.session.Domains.Create(context.Background(), "example.com")
	suite.Require().Error(error)
	suite.Require().Nil(domain)
}

func (suite *DomainErrorTestSuite) TestRead() {
	domain, error := suite.session.Domains.Read(context.Background(), "example.com")
	suite.Require().Error(error)
	suite.Require().Nil(domain)
}

func (suite *DomainErrorTestSuite) TestUpdate() {
	domain, error := suite.session.Domains.Update(context.Background(), "example.com")
	suite.Require().Error(error)
	suite.Require().Nil(domain)
}
