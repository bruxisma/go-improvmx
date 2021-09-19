package improvmx

import (
	"context"
	"strconv"

	"occult.work/doze"
)

type DomainEndpoint doze.Client

type Domain struct {
	Active            bool    `json:"active"`
	Name              string  `json:"domain"`
	Display           string  `json:"display"`
	DKIMSelector      string  `json:"dkim_selector"`
	NotificationEmail string  `json:"notification_email"`
	Whitelabel        string  `json:"white_label"`
	Added             Time    `json:"added"`
	Aliases           []Alias `json:"aliases"`
}

type DomainRecord struct {
}

// Used for creating or updating a domain entry
type DomainOption struct {
	Email string `json:"notification_email,omitempty"`
	Label string `json:"whitelabel,omitempty"`
}

// Returns a slice of domain information for the session
//
// If multiple *ListOption are passed, only the first one is used.
//
// See the API reference for more information: https://improvmx.com/api/#domains-list
func (endpoint *DomainEndpoint) List(ctx context.Context, options ...*ListOption) ([]Domain, error) {
	request := endpoint.inner().Request(ctx, &domainsResponse{})

	option := getListOption(options...)
	if error := option.validate(); error != nil {
		return nil, error
	}

	var domains []Domain
	page := 1

	if option.page != nil {
		page = *option.page
	}

	if option.startsWith != "" {
		request.SetQueryParameter("q", option.startsWith)
	}

	if option.isActive != nil {
		request.SetQueryParameter("is_active", strconv.Itoa(*option.isActive))
	}

	if option.limit != nil {
		request.SetQueryParameter("limit", strconv.Itoa(*option.limit))
	}

	for {
		request.SetQueryParameter("page", strconv.Itoa(page))
		if response, error := request.Get(domainListPath); error != nil {
			return nil, error
		} else {
			total := response.(*domainsResponse).Total
			/* Small attempt at allocation optimization */
			if len(domains) == 0 {
				domains = make([]Domain, 0, total)
			}
			domains = append(domains, response.(*domainsResponse).Domains...)
			if len(domains) >= total {
				break
			}
		}
		page++
	}
	return domains, nil
}

// Retrieve the logs for the given domain.
//
// See the API reference for more information: https://improvmx.com/api/#logs-list
func (endpoint *DomainEndpoint) Logs(ctx context.Context, domain string) ([]LogEntry, error) {
	request := endpoint.inner().Request(ctx, &logsResponse{}).
		SetPathParameter("domain", domain)
	if response, error := request.Get(domainLogsPath); error != nil {
		return nil, error
	} else {
		return response.(*logsResponse).Logs, nil
	}
}

// Adds a new domain to the given account, with an optional email and
// whitelabel. If more than one DomainOption is passed to the function, it will
// be ignored.
//
// See the API reference for more information: https://improvmx.com/api/#domains-add
func (endpoint *DomainEndpoint) Create(ctx context.Context, domain string, options ...DomainOption) (*Domain, error) {
	request := endpoint.inner().Request(ctx, &domainResponse{}).
		SetPathParameter("domain", domain).
		SetBody(getDomainOption(options...))
	if response, error := request.Post(domainCreatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*domainResponse)).Domain, nil
	}
}

func (endpoint *DomainEndpoint) Read(ctx context.Context, domain string) (*Domain, error) {
	request := endpoint.inner().Request(ctx, &domainResponse{}).
		SetPathParameter("domain", domain)
	if response, error := request.Get(domainReadPath); error != nil {
		return nil, error
	} else {
		return &(response.(*domainResponse)).Domain, nil
	}
}

// Updates the notification email or whitelabel for the given domain. If more
// than one DomainOption is passed to the function, it will be ignored
//
// It is not possible to change the domain name.
//
// See the API reference for more information: https://improvmx.com/api/#domain-update
func (endpoint *DomainEndpoint) Update(ctx context.Context, domain string, options ...DomainOption) (*Domain, error) {
	request := endpoint.inner().Request(ctx, &domainResponse{}).
		SetPathParameter("domain", domain).
		SetBody(getDomainOption(options...))
	if response, error := request.Put(domainUpdatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*domainResponse)).Domain, nil
	}
}

// Deletes the given domain from the account
//
// See the API reference for more information: https://improvmx.com/api/#domain-delete
func (endpoint *DomainEndpoint) Delete(ctx context.Context, domain string) error {
	request := endpoint.inner().Request(ctx, &deleteResponse{}).
		SetPathParameter("domain", domain)
	_, error := request.Delete(domainDeletePath)
	return error
}

// Check if the MX entries are valid for a given domain. No error is returned
// if the entries are valid.
//
// NOTE: This currently only returns an error if the domain information is not
// valid, however this is because we are swallowing the actual JSON output
// returned by the ImprovMX REST API. At some point in the future, Verify will
// be a wrapper around another function that returns the actual deserialized
// data.
func (endpoint *DomainEndpoint) Verify(ctx context.Context, domain string) error {
	request := endpoint.inner().Request(ctx, make(map[string]interface{})).
		SetPathParameter("domain", domain)
	_, error := request.Get(domainVerifyPath)
	return error
}

// getDomainOption returns either a default DomainOption *or* the first
// parameter passed in the variadic arguments
func getDomainOption(options ...DomainOption) DomainOption {
	if len(options) == 0 {
		return DomainOption{}
	}
	return options[0]
}

func (endpoint *DomainEndpoint) inner() *doze.Client {
	return (*doze.Client)(endpoint)
}
