package improvmx

import (
	"context"

	"occult.work/doze"
)

const (
	credentialsListPath   = "/domains/{domain}/credentials/"
	credentialsCreatePath = "/domains/{domain}/credentials/"
	credentialsUpdatePath = "/domains/{domain}/credentials/{username}"
	credentialsDeletePath = "/domains/{domain}/credentials/{username}"
)

type CredentialEndpoint doze.Client

type Credential struct {
	CreatedAt Time   `json:"created"`
	Username  string `json:"username"`
	Usage     int64  `json:"usage"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type credentialListResponse struct {
	credentials []Credential
}

type credentialResponse struct {
	credential Credential
}

// Returns a list of SMTP credentials for the given domain.
//
// Accessing the SMTP Credentials API is a premium feature tied to the ImprovMX
// account
//
// See the API reference for more information: https://improvmx.com/api/#smtp-list
func (endpoint *CredentialEndpoint) List(ctx context.Context, domain string) ([]Credential, error) {
	request := endpoint.inner().Request(ctx, &credentialListResponse{}).
		SetPathParameter("domain", domain)
	if response, error := request.Get(credentialsListPath); error != nil {
		return nil, error
	} else {
		return response.(*credentialListResponse).credentials, nil
	}
}

// Adds a new SMTP account to send emails
//
// See the API reference for more information: https://improvmx.com/api/#smtp-add
func (endpoint *CredentialEndpoint) Create(ctx context.Context, domain string, user User) (*Credential, error) {
	request := endpoint.inner().Request(ctx, &credentialResponse{}).
		SetPathParameter("domain", domain).
		SetBody(user)
	if response, error := request.Post(credentialsCreatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*credentialResponse)).credential, nil
	}
}

// Updates the password for the given user
//
// See the API reference for more information: https://improvmx.com/api/#smtp-update
func (endpoint *CredentialEndpoint) Update(ctx context.Context, domain string, user User) (*Credential, error) {
	request := endpoint.inner().Request(ctx, &credentialResponse{}).
		SetPathParameter("username", user.Username).
		SetPathParameter("domain", domain).
		SetBody(map[string]string{"password": user.Password})
	if response, error := request.Put(credentialsUpdatePath); error != nil {
		return nil, error
	} else {
		return &(response.(*credentialResponse)).credential, nil
	}
}

// Deletes the given SMTP user for the given domain
//
// See the API reference for more information: https://improvmx.com/api/#smtp-delete
func (endpoint *CredentialEndpoint) Delete(ctx context.Context, domain, username string) error {
	request := endpoint.inner().Request(ctx, &deleteResponse{}).
		SetPathParameter("username", username).
		SetPathParameter("domain", domain)
	_, error := request.Delete(credentialsDeletePath)
	return error
}

func (endpoint *CredentialEndpoint) inner() *doze.Client {
	return (*doze.Client)(endpoint)
}
