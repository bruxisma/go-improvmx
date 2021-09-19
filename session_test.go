package improvmx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"occult.work/doze"
)

func TestSessionConstruction(t *testing.T) {
	New("token", WithBaseURL("https://fake-base-url"))
	New("token", WithUserAgent("Agent"))
	New("token", WithClient(nil))
}

func TestSessionWithDebug(test *testing.T) {
	assert := assert.New(test)
	session, error := New("token", WithDebug())
	assert.NoError(error)
	assert.NotEmpty(session)
	assert.NotEmpty(session.Credentials)
	assert.NotEmpty(session.Account)
	assert.NotEmpty(session.Domains)
	assert.NotEmpty(session.Aliases)
}

func TestSessionWithUserAgent(test *testing.T) {
	assert := assert.New(test)
	session, error := New("token", WithUserAgent("Agent-Name"))
	assert.NoError(error)
	assert.NotEmpty(session)
	assert.NotEmpty(session.Credentials)
	assert.NotEmpty(session.Account)
	assert.NotEmpty(session.Domains)
	assert.NotEmpty(session.Aliases)
}

func TestSessionWithNilClient(test *testing.T) {
	assert := assert.New(test)
	session, error := New("token", WithClient(nil))
	assert.Error(error)
	assert.Empty(session)
}

func TestSessionWithClient(test *testing.T) {
	assert := assert.New(test)
	client := doze.NewClient()
	session, error := New("token", WithClient(client))
	assert.NoError(error)
	assert.NotEmpty(session)
	assert.NotEmpty(session.Credentials)
	assert.NotEmpty(session.Account)
	assert.NotEmpty(session.Domains)
	assert.NotEmpty(session.Aliases)
}
