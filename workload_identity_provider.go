package dtrack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WorkloadIdentityProviderType string

const (
	WorkloadIdentityProviderTypeOIDC   WorkloadIdentityProviderType = "OIDC"
	WorkloadIdentityProviderTypeSPIFFE WorkloadIdentityProviderType = "SPIFFE"
)

type CreateWorkloadIdentityProviderRequest struct {
	Name                   string                       `json:"name"`
	Type                   WorkloadIdentityProviderType `json:"type"`
	Issuer                 string                       `json:"issuer"`
	Audience               string                       `json:"audience"`
	JWKSURL                string                       `json:"jwks_url,omitempty"`
	JWKS                   json.RawMessage              `json:"jwks,omitempty"`
	SessionLifetimeSeconds *int32                       `json:"session_lifetime_seconds,omitempty"`
}

type CreateWorkloadIdentityProviderResponse struct {
	URL string // Returned within `Location` header, hence why no json tag.
}

type WorkloadIdentityProvider struct {
	Name                   string                       `json:"name"`
	Type                   WorkloadIdentityProviderType `json:"type"`
	Issuer                 string                       `json:"issuer"`
	Audience               string                       `json:"audience"`
	JWKSURL                string                       `json:"jwks_url,omitempty"`
	JWKSKeyIDs             []string                     `json:"jwks_key_ids,omitempty"`
	SessionLifetimeSeconds int32                        `json:"session_lifetime_seconds"`
	CreatedAt              Timestamp                    `json:"created_at"`
}

type UpdateWorkloadIdentityProviderRequest struct {
	Issuer                 *string         `json:"issuer,omitempty"`
	Audience               *string         `json:"audience,omitempty"`
	JWKSURL                *string         `json:"jwks_url,omitempty"`
	JWKS                   json.RawMessage `json:"jwks,omitempty"`
	SessionLifetimeSeconds *int32          `json:"session_lifetime_seconds,omitempty"`
}

type WorkloadIdentityProviderService struct {
	client *Client
}

func (wips WorkloadIdentityProviderService) Create(ctx context.Context, provider CreateWorkloadIdentityProviderRequest) (r CreateWorkloadIdentityProviderResponse, err error) {
	err = wips.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := wips.client.newRequest(ctx, http.MethodPost, "api/v2/workload-identity-providers", withBody(provider), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	res, err := wips.client.doRequest(req, nil)
	if err != nil {
		return
	}

	r.URL = res.Header.Get("Location")
	return
}

func (wips WorkloadIdentityProviderService) List(ctx context.Context, query, pageToken string, limit uint32) (p PageV5[WorkloadIdentityProvider], err error) {
	err = wips.client.assertServerVersionAtLeast("5.2.0")
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

	req, err := wips.client.newRequest(ctx, http.MethodGet, "api/v2/workload-identity-providers", withParams(params))
	if err != nil {
		return
	}

	_, err = wips.client.doRequest(req, &p)
	return
}

func (wips WorkloadIdentityProviderService) Get(ctx context.Context, name string) (r WorkloadIdentityProvider, err error) {
	err = wips.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := wips.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("api/v2/workload-identity-providers/%s", url.PathEscape(name)))
	if err != nil {
		return
	}

	_, err = wips.client.doRequest(req, &r)
	return
}

func (wips WorkloadIdentityProviderService) Update(ctx context.Context, name string, provider UpdateWorkloadIdentityProviderRequest) (err error) {
	err = wips.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := wips.client.newRequest(ctx, http.MethodPatch, fmt.Sprintf("api/v2/workload-identity-providers/%s", url.PathEscape(name)), withBody(provider), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	_, err = wips.client.doRequest(req, nil)
	return
}

func (wips WorkloadIdentityProviderService) Delete(ctx context.Context, name string) (err error) {
	err = wips.client.assertServerVersionAtLeast("5.2.0")
	if err != nil {
		return
	}

	req, err := wips.client.newRequest(ctx, http.MethodDelete, fmt.Sprintf("api/v2/workload-identity-providers/%s", url.PathEscape(name)), withAcceptContentType("*/*"))
	if err != nil {
		return
	}

	_, err = wips.client.doRequest(req, nil)
	return
}
