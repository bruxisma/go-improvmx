package internal

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

/* This is used specifically for tests */
type Suite struct {
	suite.Suite
	Server *httptest.Server
	Data   *embed.FS
}

func (suite *Suite) TearDownSuite() {
	if suite.Server != nil {
		suite.Server.Close()
	}
}

// We only support GET requests at the moment. This needs to change soon to
// support POST, PATCH, DELETE, etc.
func (suite *Suite) Initialize(testData map[string][]byte) {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		suite.Require().Fail("Unexpected request URL for request")
	})
	for path, slice := range testData {
		/* If this is not done, we get strange and spurious failures */
		data := slice[:]
		router.Path(path).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			writer.Write(data)
		})
	}
	suite.Server = httptest.NewServer(router)
}

func (suite *Suite) InitializeWithRouter(router *Router) {
	router.inner().HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		suite.Require().Fail("Unexpected path requested")
	})
	suite.Server = httptest.NewServer((*mux.Router)(router))
}

func (suite *Suite) Skip() {
	suite.T().SkipNow()
}

func (suite *Suite) FileResponseHandler(filename string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, error := suite.Data.ReadFile(filename)
		suite.Require().NoError(error)
		writer.Write(data)
	}
}

func (suite *Suite) Template(patterns ...string) *template.Template {
	template, error := template.ParseFS(suite.Data, patterns...)
	suite.Require().NoError(error)
	return template
}

func (suite *Suite) Render(writer io.Writer, path string, data interface{}) {
	template, error := template.ParseFS(suite.Data, path)
	suite.Require().NoError(error)
	error = template.Execute(writer, data)
	suite.Require().NoError(error)
}

func (suite *Suite) Parameters(request *http.Request) map[string]string {
	var parameters map[string]string
	decoder := json.NewDecoder(request.Body)
	suite.Require().NoError(decoder.Decode(&parameters))
	return parameters
}
