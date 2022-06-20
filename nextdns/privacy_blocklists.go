package nextdns

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// privacyBlocklistsAPIPath is the HTTP path for the privacy blocklist API.
const privacyBlocklistsAPIPath = "privacy/blocklists"

// PrivacyBlocklists represents a privacy blocklist of a profile.
type PrivacyBlocklists struct {
	ID        string     `json:"id"`
	Name      string     `json:"name,omitempty"`
	Website   string     `json:"website,omitempty"`
	Entries   int        `json:"entries,omitempty"`
	UpdatedOn *time.Time `json:"updatedOn,omitempty"`
}

// CreatePrivacyBlocklistsRequest encapsulates the request for creating a privacy blocklist.
type CreatePrivacyBlocklistsRequest struct {
	Profile string
}

// GetPrivacyBlocklistsRequest encapsulates the request for getting the privacy blocklist.
type GetPrivacyBlocklistsRequest struct {
	Profile string
}

// PrivacyBlocklistsService is an interface for communicating with the NextDNS privacy blocklist API endpoint.
type PrivacyBlocklistsService interface {
	Create(context.Context, *CreatePrivacyBlocklistsRequest, []*PrivacyBlocklists) error
	Get(context.Context, *GetPrivacyBlocklistsRequest) ([]*PrivacyBlocklists, error)
}

// privacyBlocklistsResponse represents the NextDNS privacy blocklist service.
type privacyBlocklistsResponse struct {
	PrivacyBlocklists []*PrivacyBlocklists `json:"data"`
}

// privacyBlocklistsService represents the NextDNS privacy blocklist service.
type privacyBlocklistsService struct {
	client *Client
}

var _ PrivacyBlocklistsService = &privacyBlocklistsService{}

// NewPrivacyBlocklistsService returns a new NextDNS privacy blocklist service.
// nolint: revive
func NewPrivacyBlocklistsService(client *Client) *privacyBlocklistsService {
	return &privacyBlocklistsService{
		client: client,
	}
}

// Create creates a privacy blocklist list for a profile.
func (s *privacyBlocklistsService) Create(ctx context.Context, request *CreatePrivacyBlocklistsRequest, v []*PrivacyBlocklists) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), privacyBlocklistsAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to create a privacy blocklist")
	}

	response := privacyBlocklistsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to create a privacy blocklist")
	}

	return nil
}

// Get returns the privacy blocklist for a profile.
func (s *privacyBlocklistsService) Get(ctx context.Context, request *GetPrivacyBlocklistsRequest) ([]*PrivacyBlocklists, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), privacyBlocklistsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the privacy blocklist")
	}

	response := privacyBlocklistsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the privacy blocklist")
	}

	return response.PrivacyBlocklists, nil
}
