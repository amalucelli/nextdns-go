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
	ProfileID         string
	PrivacyBlocklists []*PrivacyBlocklists
}

// ListPrivacyBlocklistsRequest encapsulates the request for getting the privacy blocklist.
type ListPrivacyBlocklistsRequest struct {
	ProfileID string
}

// PrivacyBlocklistsService is an interface for communicating with the NextDNS privacy blocklist API endpoint.
type PrivacyBlocklistsService interface {
	Create(context.Context, *CreatePrivacyBlocklistsRequest) error
	List(context.Context, *ListPrivacyBlocklistsRequest) ([]*PrivacyBlocklists, error)
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
func (s *privacyBlocklistsService) Create(ctx context.Context, request *CreatePrivacyBlocklistsRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyBlocklistsAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.PrivacyBlocklists)
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

// List returns the privacy blocklist for a profile.
func (s *privacyBlocklistsService) List(ctx context.Context, request *ListPrivacyBlocklistsRequest) ([]*PrivacyBlocklists, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyBlocklistsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to list the privacy blocklist")
	}

	response := privacyBlocklistsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to list the privacy blocklist")
	}

	return response.PrivacyBlocklists, nil
}
