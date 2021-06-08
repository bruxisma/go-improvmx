package improvmx

import (
	"errors"
	"fmt"

	"occult.work/improvmx/internal"
)

const (
	// The email was accepted to be processed
	Queued MessageStatus = "QUEUED"

	// The email was refused at the SMTP connection
	Refused = "REFUSED"

	// The email was successfully delivered to the end destination
	Delievered = "DELIVERED"

	// The end destination refused the email temporarily. ImprovMX will try again
	// multiple times with increased delay between
	SoftBounce = "SOFT-BOUNCE"

	// The end destination couldn't accept the email definitively
	HardBounce = "HARD-BOUNCE"
)

const (
	accountLabelsPath = "/account/whitelabels/"
	accountReadPath   = "/account/"

	aliasListPath   = "/domains/{domain}/aliases/"
	aliasCreatePath = "/domains/{domain}/aliases/"
	aliasReadPath   = "/domains/{domain}/aliases/{alias}/"
	aliasUpdatePath = "/domains/{domain}/aliases/{alias}/"
	aliasDeletePath = "/domains/{domain}/aliases/{alias}/"
	aliasLogsPath   = "/domains/{domain}/logs/{alias}/"

	credentialsListPath   = "/domains/{domain}/credentials/"
	credentialsCreatePath = "/domains/{domain}/credentials/"
	credentialsUpdatePath = "/domains/{domain}/credentials/{username}/"
	credentialsDeletePath = "/domains/{domain}/credentials/{username}/"

	domainListPath   = "/domains/"
	domainCreatePath = "/domains/{domain}/"
	domainReadPath   = "/domains/{domain}/"
	domainUpdatePath = "/domains/{domain}/"
	domainDeletePath = "/domains/{domain}/"
	domainVerifyPath = "/domains/{domain}/check/"
	domainLogsPath   = "/domains/{domain}/logs/"
)

type Time = internal.Time
type MessageStatus string

type Contact struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Error struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

// Used by the Domain and Alias List methods to filter additional results.
//
// Due to the nature of the ImprovMX REST API, users must use the provided
// methods to set the ListOption. However, the Alias.List method might still
// error if the SetLimit method has been called.
type ListOption struct {
	startsWith string
	isActive   *int
	limit      *int
	page       *int
}

type LogEvent struct {
	Code      int64         `json:"code"`
	CreatedAt string        `json:"created"`
	ID        string        `json:"id"`
	Local     string        `json:"local"`
	Message   string        `json:"message"`
	Server    string        `json:"server"`
	Status    MessageStatus `json:"status"`
}

type LogEntry struct {
	CreatedAt  string     `json:"created"`
	CreatedRaw string     `json:"created_raw"`
	Events     []LogEvent `json:"events"`
	Address    Contact    `json:"forward"`
	Hostname   string     `json:"hostname"`
	ID         string     `json:"id"`
	MessageID  string     `json:"messageId"`
	Recipient  Contact    `json:"recipient"`
	Sender     Contact    `json:"sender"`
	Subject    string     `json:"subject"`
	Transport  string     `json:"transport"`
}
type accountResponse struct {
	Account Account
	Success bool
}

type whitelabelsResponse struct {
	Whitelabels []Whitelabel
	Success     bool
}

type domainsResponse struct {
	Domains []Domain
	Success bool
	Total   int
	Limit   int
	Page    int
}

type domainResponse struct {
	Domain  Domain
	Success bool
}

type aliasesResponse struct {
	Aliases []Alias
	Success bool
	Total   int
	Page    int
}

type aliasResponse struct {
	Alias   Alias
	Success bool
}

type deleteResponse struct {
	Success bool
}

type logsResponse struct {
	Success bool
	Logs    []LogEntry
}

func NewListOption() *ListOption {
	return &ListOption{"", nil, nil, nil}
}

// Clears the value provided by SetStartsWith
func (option *ListOption) ClearStartsWith() *ListOption {
	option.startsWith = ""
	return option
}

// Clears the value provided by SetIsActive
func (option *ListOption) ClearIsActive() *ListOption {
	option.isActive = nil
	return option
}

// Clears the value provided by SetLimit
func (option *ListOption) ClearLimit() *ListOption {
	option.limit = nil
	return option
}

// Clears the value provided by SetPage
func (option *ListOption) ClearPage() *ListOption {
	option.page = nil
	return option
}

// Used to limit searches to domains or aliases that start with the provided
// string
func (option *ListOption) SetStartsWith(value string) *ListOption {
	option.startsWith = value
	return option
}

// Used to return only active or inactive domains. If unset, List methods will
// return all domains.
func (option *ListOption) SetIsActive(value bool) *ListOption {
	active := 0
	if value {
		active = 1
	}
	option.isActive = &active
	return option
}

// Number of domains to return. Defaults to 50, must be between 5 and 100.  If
// this value is set when calling the Alias.List method, the method will fail
// with an error
func (option *ListOption) SetLimit(value int) *ListOption {
	option.limit = &value
	return option
}

// Current page to load. Must be greater or equal to 1, otherwise List methods
// will fail with an error
func (option *ListOption) SetPage(value int) *ListOption {
	option.page = &value
	return option
}

func (option *ListOption) validate() error {
	if option.limit != nil {
		limit := *option.limit
		if limit < 5 || limit > 100 {
			return errors.New(fmt.Sprintf("Limit is outside the expected range of [5, 100]: %d", limit))
		}
	}

	if option.page != nil {
		page := *option.page
		if page < 1 {
			return errors.New(fmt.Sprintf("Page must be greater than or equal to 1: %d", page))
		}
	}

	return nil
}

func getListOption(options ...ListOption) *ListOption {
	if len(options) == 0 {
		return nil
	}
	return &options[0]
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
