package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// denylistAPIPath is the HTTP path for the denylist API.
const denylistAPIPath = "denylist"

// Denylist represents the denylist of a profile.
type Denylist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active"`
}

// CreateDenylistRequest encapsulates the request for creating a denylist.
type CreateDenylistRequest struct {
	ProfileID string
	Denylist  []*Denylist
}

// GetDenylistRequest encapsulates the request for getting a denylist.
type GetDenylistRequest struct {
	ProfileID string
}

// UpdateDenylistRequest encapsulates the request for updating a denylist.
type UpdateDenylistRequest struct {
	ProfileID string
	ID        string
	Denylist  *Denylist
}

// DenylistService is an interface for communicating with the NextDNS denylist API endpoint.
type DenylistService interface {
	Create(context.Context, *CreateDenylistRequest) error
	Get(context.Context, *GetDenylistRequest) ([]*Denylist, error)
	Update(context.Context, *UpdateDenylistRequest) error
}

// denylistResponse represents the denylist response.
type denylistResponse struct {
	Denylist []*Denylist `json:"data"`
}

// denylistService represents the NextDNS denylist service.
type denylistService struct {
	client *Client
}

var _ DenylistService = &denylistService{}

// NewDenylistService returns a new NextDNS denylist service.
// nolint: revive
func NewDenylistService(client *Client) *denylistService {
	return &denylistService{
		client: client,
	}
}

// Create creates a denylist for a profile.
func (s *denylistService) Create(ctx context.Context, request *CreateDenylistRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), denylistAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.Denylist)
	if err != nil {
		return errors.Wrap(err, "error creating request to create an deny list")
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "error making a request to create an deny list")
	}

	return nil
}

// Get returns the denylist of a profile.
func (s *denylistService) Get(ctx context.Context, request *GetDenylistRequest) ([]*Denylist, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), denylistAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the deny list")
	}

	response := denylistResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the deny list")
	}

	return response.Denylist, nil
}

// Update updates a denylist of a profile.
func (s *denylistService) Update(ctx context.Context, request *UpdateDenylistRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), denylistIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.Denylist)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error creating request to update the deny list id: %s", request.ID))
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error making a request to update the deny list id: %s", request.ID))
	}

	return nil
}

// denylistIDAPIPath returns the HTTP path for the denylist API.
func denylistIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", denylistAPIPath, id)
}
