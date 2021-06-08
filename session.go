package improvmx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"occult.work/improvmx/internal"
)

const (
	domainsPath = "/domains/"
	domainPath  = "/domains/{domain}"
	//BaseURLv3 is the base URL for the improvmx API
	BaseURLv3 = "https://api.improvmx.com/v3"
)

type Session struct {
	client      *resty.Client
	Credentials *CredentialEndpoint
	Account     *AccountEndpoint
	Domains     *DomainEndpoint
	Aliases     *AliasEndpoint
}

type SessionOption func(*Session) error

func (session *Session) request(ctx context.Context, result interface{}) *internal.Request {
	r := session.client.
		NewRequest().
		ForceContentType("application/json").
		SetContext(ctx).
		SetResult(result)
	return (*internal.Request)(r)
}

func New(token string, options ...SessionOption) (*Session, error) {
	session := &Session{}
	session.client = resty.New().
		SetTimeout(time.Duration(1 * time.Minute)).
		SetAuthToken(fmt.Sprintf("api:%s", token)).
		SetAuthScheme("Basic").
		SetHostURL(BaseURLv3).
		SetError(&Error{})
	for _, option := range options {
		if error := option(session); error != nil {
			return nil, error
		}
	}
	session.Credentials = (*CredentialEndpoint)(session.client)
	session.Account = (*AccountEndpoint)(session.client)
	session.Domains = (*DomainEndpoint)(session.client)
	session.Aliases = (*AliasEndpoint)(session.client)
	return session, nil
}

// Sets the session client directly. This is provided for extreme user needs.
// However, the resty.Client.SetError type is still set to Error, and the host
// url is set toe BaseURLv3. Use WithHostURL after calling WithClient if
// overriding the Host URL is required.
func WithClient(client *resty.Client) SessionOption {
	return func(session *Session) error {
		client.SetError(&Error{}).SetHostURL(BaseURLv3)
		session.client = client
		return nil
	}
}

// Sets the User-Agent to the given string on the session. By default, this is
// unset.
func WithUserAgent(agent string) SessionOption {
	return func(session *Session) error {
		session.client.SetHeader("User-Agent", agent)
		return nil
	}
}

// Sets the Base URL to query. This is typically BaseURLv3, however for both
// testing and possible future sandboxing support, this option is provided.
func WithHostURL(url string) SessionOption {
	return func(session *Session) error {
		session.client.SetHostURL(url)
		return nil
	}
}

// Enables the debug mode on the Session's internal http client
func WithDebug() SessionOption {
	return func(session *Session) error {
		session.client.SetDebug(true)
		return nil
	}
}
