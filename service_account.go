package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type CreateServiceAccountRequest struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
}

type CreateServiceAccountResponse struct {
	URL string // Returned within `Location` header, hence why no json tag.
}

type ServiceAccountTeam struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type ServiceAccount struct {
	Name        string               `json:"name"`
	Username    string               `json:"username"`
	Email       string               `json:"email,omitempty"`
	Suspended   bool                 `json:"suspended"`
	Teams       []ServiceAccountTeam `json:"teams,omitempty"`
	Permissions []string             `json:"permissions,omitempty"`
}

type UpdateServiceAccountRequest struct {
	Email     *string `json:"email,omitempty"`
	Suspended *bool   `json:"suspended,omitempty"`
}

type WorkloadIdentityBinding struct {
	UUID         uuid.UUID  `json:"uuid"`
	ProviderName string     `json:"provider_name"`
	Subject      string     `json:"subject"`
	Condition    string     `json:"condition,omitempty"`
	CreatedAt    Timestamp  `json:"created_at"`
	LastUsedAt   *Timestamp `json:"last_used_at,omitempty"`
}

type CreateWorkloadIdentityBindingRequest struct {
	ProviderName string `json:"provider_name"`
	Subject      string `json:"subject"`
	Condition    string `json:"condition,omitempty"`
}

type CreateWorkloadIdentityBindingResponse struct {
	URL  string    // Returned within `Location` header, hence why no json tag.
	UUID uuid.UUID `json:"uuid"`
}

type ServiceAccountService struct {
	client *Client
}

func (sas ServiceAccountService) Create(ctx context.Context, account CreateServiceAccountRequest) (r CreateServiceAccountResponse, err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodPost, "api/v2/service-accounts", withBody(account), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	res, err := sas.client.doRequest(req, nil)
	if err != nil {
		return
	}

	r.URL = res.Header.Get("Location")
	return
}

func (sas ServiceAccountService) List(ctx context.Context, query, pageToken string, limit uint32) (p PageV5[ServiceAccount], err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	// Empty page tokens are rejected by the server, so are only sent when present.
	params := map[string]string{
		"limit": fmt.Sprintf("%d", limit),
	}
	if pageToken != "" {
		params["page_token"] = pageToken
	}
	if query != "" {
		params["q"] = query
	}

	req, err := sas.client.newRequest(ctx, http.MethodGet, "api/v2/service-accounts", withParams(params))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, &p)
	return
}

func (sas ServiceAccountService) Get(ctx context.Context, name string) (r ServiceAccount, err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/service-accounts/%s", url.PathEscape(name)))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, &r)
	return
}

func (sas ServiceAccountService) Update(ctx context.Context, name string, account UpdateServiceAccountRequest) (err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodPatch, fmt.Sprintf("api/v2/service-accounts/%s", url.PathEscape(name)), withBody(account), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, nil)
	return
}

func (sas ServiceAccountService) Delete(ctx context.Context, name string) (err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v2/service-accounts/%s", url.PathEscape(name)), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, nil)
	return
}

func (sas ServiceAccountService) CreateWorkloadIdentityBinding(ctx context.Context, name string, binding CreateWorkloadIdentityBindingRequest) (r CreateWorkloadIdentityBindingResponse, err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodPost, fmt.Sprintf("api/v2/service-accounts/%s/workload-identity-bindings", url.PathEscape(name)), withBody(binding))
	if err != nil {
		return
	}

	res, err := sas.client.doRequest(req, &r)
	if err != nil {
		return
	}

	r.URL = res.Header.Get("Location")
	return
}

func (sas ServiceAccountService) ListWorkloadIdentityBindings(ctx context.Context, name, pageToken string, limit uint32) (p PageV5[WorkloadIdentityBinding], err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	// Empty page tokens are rejected by the server, so are only sent when present.
	params := map[string]string{
		"limit": fmt.Sprintf("%d", limit),
	}
	if pageToken != "" {
		params["page_token"] = pageToken
	}

	req, err := sas.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/service-accounts/%s/workload-identity-bindings", url.PathEscape(name)), withParams(params))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, &p)
	return
}

func (sas ServiceAccountService) DeleteWorkloadIdentityBinding(ctx context.Context, name string, bindingID uuid.UUID) (err error) {
	err = sas.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := sas.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v2/service-accounts/%s/workload-identity-bindings/%s", url.PathEscape(name), bindingID), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	_, err = sas.client.doRequest(req, nil)
	return
}
