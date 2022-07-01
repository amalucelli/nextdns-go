package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// allowlistAPIPath is the HTTP path for the allowlist API.
const allowlistAPIPath = "allowlist"

// Allowlist represents the allow list of a profile.
type Allowlist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active"`
}

// CreateAllowlistRequest encapsulates the request for creating an allowlist.
type CreateAllowlistRequest struct {
	ProfileID string
	Allowlist []*Allowlist
}

// GetAllowlistRequest encapsulates the request for getting an allowlist.
type GetAllowlistRequest struct {
	ProfileID string
}

// UpdateAllowlistRequest encapsulates the request for updating an allowlist.
type UpdateAllowlistRequest struct {
	ProfileID string
	ID        string
	Allowlist *Allowlist
}

// AllowlistService is an interface for communicating with the NextDNS allowlist API endpoint.
type AllowlistService interface {
	Create(context.Context, *CreateAllowlistRequest) error
	Get(context.Context, *GetAllowlistRequest) ([]*Allowlist, error)
	Update(context.Context, *UpdateAllowlistRequest) error
}

// allowlistResponse represents the allowlist response.
type allowlistResponse struct {
	Allowlist []*Allowlist `json:"data"`
}

// privacyService represents the NextDNS allowlist service.
type allowlistService struct {
	client *Client
}

var _ AllowlistService = &allowlistService{}

// NewAllowlistService returns a new NextDNS allowlist service.
// nolint: revive
func NewAllowlistService(client *Client) *allowlistService {
	return &allowlistService{
		client: client,
	}
}

// Create creates an allowlist for a profile.
func (s *allowlistService) Create(ctx context.Context, request *CreateAllowlistRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), allowlistAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.Allowlist)
	if err != nil {
		return errors.Wrap(err, "error creating request to create an allow list")
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "error making a request to create an allow list")
	}

	return nil
}

// Get returns the allowlist of a profile.
func (s *allowlistService) Get(ctx context.Context, request *GetAllowlistRequest) ([]*Allowlist, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), allowlistAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the allow list")
	}

	response := allowlistResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the allow list")
	}

	return response.Allowlist, nil
}

// Update updates an allowlist of a profile.
func (s *allowlistService) Update(ctx context.Context, request *UpdateAllowlistRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), allowlistIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.Allowlist)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error creating request to update the allow list id: %s", request.ID))
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error making a request to update the allow list id: %s", request.ID))
	}

	return nil
}

// allowlistIDAPIPath returns the HTTP path for the allowlist API.
func allowlistIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", allowlistAPIPath, id)
}
