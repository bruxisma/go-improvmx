package improvmx

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"testing"

	"occult.work/doze/test"
)

type AliasTestSuite struct {
	test.Suite
	session *Session
}

const createAliasResponse = `
	{
		"alias": {
			"forward": "{{ .Address }}",
			"alias": "{{ .Name }}",
			"id": {{ .ID }}
		},
		"success": true
	}
`

/* This is all *very* hacky, and we cannot (unfortunately) easily throw
* parameterized tests at this problem
 */
func (suite *AliasTestSuite) SetupSuite() {
	router := test.NewRouter().
		Get(aliasReadPath, suite.FileResponseHandler("testdata/alias/read.json")).
		Get(aliasLogsPath, suite.FileResponseHandler("testdata/alias/logs.json")).
		Get(aliasListPath, suite.FileResponseHandler("testdata/alias/list.json")).
		Post(aliasCreatePath, func(writer http.ResponseWriter, request *http.Request) {
			parameters := suite.Parameters(request)
			suite.Render(writer, "testdata/alias/create.json.tmpl", Alias{
				parameters["forward"],
				parameters["alias"],
				47,
			})
		}).
		Delete(aliasDeletePath, func(writer http.ResponseWriter, request *http.Request) {
			suite.Require().Equal(request.URL.Path, "/domains/example.com/aliases/richard/")
			fmt.Fprintf(writer, `{ "success": true }`)
		})
	suite.Initialize(router)
	suite.Data = &testData
	suite.session = setupSession(suite.Server)
}

func TestAlias(t *testing.T) {
	test.Run(t, new(AliasTestSuite))
}

func (suite *AliasTestSuite) TestList() {
	aliases, error := suite.session.Aliases.List(context.Background(), "example.com")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(aliases)
}

func (suite *AliasTestSuite) TestLogs() {
	logs, error := suite.session.Aliases.Logs(context.Background(), "example.com", "richard")
	suite.Require().NoError(error)
	suite.Require().NotEmpty(logs)
}

func (suite *AliasTestSuite) TestCreate() {
	alias, error := suite.session.Aliases.Create(context.Background(), "example.com", "richard", "richard@example.test")
	suite.Require().NoError(error)
	suite.Equal("richard@example.test", alias.Address)
	suite.Equal("richard", alias.Name)
}

func (suite *AliasTestSuite) TestRead() {
	alias, error := suite.session.Aliases.Read(context.Background(), "example.com", "richard")
	suite.Require().NoError(error)
	suite.Equal("richard.hendricks@example.com", alias.Address)
	suite.Equal("richard", alias.Name)
}

func (suite *AliasTestSuite) TestUpdate() {
	suite.Skip()
}

func (suite *AliasTestSuite) TestDelete() {
	error := suite.session.Aliases.Delete(context.Background(), "example.com", "richard")
	suite.Require().NoError(error)
}
