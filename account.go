package improvmx

import (
	"context"

	"occult.work/doze"
)

type AccountEndpoint doze.Client

type AccountPlan struct {
	AliasesLimit int64  `json:"aliases_limit"`
	DailyQuota   int64  `json:"daily_quota"`
	DomainsLimit int64  `json:"domains_limit"`
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	Yearly       bool   `json:"yearly"`
}

type Account struct {
	BillingEmail   string `json:"billing_email"`
	CancelsOn      Time   `json:"cancels_on"`
	CardBrand      string `json:"card_brand"`
	CompanyDetails string `json:"company_details"`
	CompanyName    string `json:"company_name"`
	CompanyVAT     string `json:"company_vat"`
	Country        string `json:"country"`
	CreatedAt      Time   `json:"created"`
	Email          string `json:"email"`
	Last4          string `json:"last4"`
	Limits         struct {
		Aliases      int `json:"aliases"`
		DailyQuota   int `json:"daily_quota"`
		Domains      int `json:"domains"`
		RateLimit    int `json:"ratelimit"`
		Redirections int `json:"redirections"`
		Subdomains   int `json:"subdomains"`
	} `json:"limits"`
	LockReason   string       `json:"lock_reason"`
	Locked       bool         `json:"locked"`
	Password     bool         `json:"password"`
	Plan         *AccountPlan `json:"plan"`
	Premium      bool         `json:"premium"`
	PrivacyLevel int64        `json:"privacy_level"`
	RenewDate    Time         `json:"renew_date"`
}

type Whitelabel struct {
	Name string `json:"name"`
}

// Returns the domains used as whitelabels associated with the account.
func (endpoint *AccountEndpoint) Labels(ctx context.Context) ([]Whitelabel, error) {
	request := endpoint.inner().Request(ctx, &whitelabelsResponse{
		Whitelabels: make([]Whitelabel, 0),
	})
	if response, error := request.Get(accountLabelsPath); error != nil {
		return nil, error
	} else {
		return response.(*whitelabelsResponse).Whitelabels, nil
	}
}

// Returns an instance of the Account.
func (endpoint *AccountEndpoint) Read(ctx context.Context) (*Account, error) {
	request := endpoint.inner().Request(ctx, new(accountResponse))
	if response, error := request.Get(accountReadPath); error != nil {
		return nil, error
	} else {
		return &(response.(*accountResponse)).Account, nil
	}
}

func (endpoint *AccountEndpoint) inner() *doze.Client {
	return (*doze.Client)(endpoint)
}
