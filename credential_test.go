package improvmx

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/suite"
	"occult.work/improvmx/internal"
)

type CredentialTestSuite struct {
	internal.Suite
	session *Session
}

func (suite *CredentialTestSuite) SetupSuite() {
	suite.Initialize(make(map[string][]byte))
	suite.session, _ = New("credential-test-suite", WithHostURL(suite.Server.URL))
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
	suite.Skip()
}
