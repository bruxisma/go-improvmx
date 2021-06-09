package improvmx

import (
	"context"
	"strconv"

	"occult.work/improvmx/internal"
)

// Used to interact with the email alias specific parts of the ImprovMX REST
// API
type AliasEndpoint internal.Client

type Alias struct {
	// The address the Alias will forward to
	Address string `json:"forward"`
	Name    string `json:"alias"`
	ID      int64  `json:"id"`
}

// Returns a slice of Alias for the given domain.
//
// If multiple *ListOption are passed, only the first one is used.
//
// See the API reference for more information: https://improvmx.com/api/#alias-list
func (endpoint *AliasEndpoint) List(ctx context.Context, domain string, options ...*ListOption) ([]Alias, error) {
	request := endpoint.inner().Request(ctx, &aliasesResponse{}).
		SetPathParameter("domain", domain)

	option := getListOption(options...)
	if error := option.validate(); error != nil {
		return nil, error
	}

	var aliases []Alias
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

	for {
		request.SetQueryParameter("page", strconv.Itoa(page))
		if response, error := request.Get(aliasListPath); error != nil {
			return nil, error
		} else {
			aliases = append(aliases, response.(*aliasesResponse).Aliases...)
			if len(aliases) >= response.(*aliasesResponse).Total {
				break
			}
		}
		page++
	}
	return aliases, nil
}

// Retrieve the logs for the given alias in the given domain.
//
// See the API reference for more information:
// https://improvmx.com/api/#logs-alias
func (endpoint *AliasEndpoint) Logs(ctx context.Context, domain, alias string) ([]LogEntry, error) {
	request := endpoint.inner().Request(ctx, &logsResponse{}).
		SetPathParameter("domain", domain).
		SetPathParameter("alias", alias)
	if response, error := request.Get(aliasLogsPath); error != nil {
		return nil, error
	} else {
		return response.(*logsResponse).Logs, nil
	}
}

// Creates a new alias for the given domain
//
// See the API reference for more information:
// https://improvmx.com/api/#alias-add
func (endpoint *AliasEndpoint) Create(ctx context.Context, domain, alias, address string) (*Alias, error) {
	request := endpoint.inner().Request(ctx, &aliasResponse{}).
		SetPathParameter("domain", domain).
		SetBody(map[string]string{"alias": alias, "forward": address})
	if response, error := request.Post(aliasCreatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*aliasResponse)).Alias, nil
	}
}

// Returns a single alias object for the named domain and alias
//
// See the API reference for more information:
// https://improvmx.com/api/#alias-details
func (endpoint *AliasEndpoint) Read(ctx context.Context, domain, alias string) (*Alias, error) {
	request := endpoint.inner().Request(ctx, &aliasResponse{}).
		SetPathParameter("domain", domain).
		SetPathParameter("alias", alias)
	if response, error := request.Get(aliasReadPath); error != nil {
		return nil, error
	} else {
		return &(response.(*aliasResponse)).Alias, nil
	}
}

// Updates a single alias for the given domain and alias
//
// See the API reference for more information:
// https://improvmx.com/api/#alias-update
func (endpoint *AliasEndpoint) Update(ctx context.Context, domain, alias, address string) (*Alias, error) {
	request := endpoint.inner().Request(ctx, &aliasResponse{}).
		SetPathParameter("domain", domain).
		SetPathParameter("alias", alias).
		SetBody(map[string]string{"forward": address})
	if response, error := request.Put(aliasUpdatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*aliasResponse)).Alias, nil
	}
}

// Deletes a single for the given domain and alias
//
// See the API reference for more information:
// https://improvmx.com/api/#alias-delete
func (endpoint *AliasEndpoint) Delete(ctx context.Context, domain, alias string) error {
	request := endpoint.inner().Request(ctx, &deleteResponse{}).
		SetPathParameter("domain", domain).
		SetPathParameter("alias", alias)
	if _, error := request.Delete(aliasDeletePath); error != nil {
		return error
	}
	return nil
}

func (endpoint *AliasEndpoint) inner() *internal.Client {
	return (*internal.Client)(endpoint)
}
