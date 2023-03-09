package nextdns

import (
	"context"
	"fmt"
	"net/http"
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

// ListDenylistRequest encapsulates the request for getting a denylist.
type ListDenylistRequest struct {
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
	List(context.Context, *ListDenylistRequest) ([]*Denylist, error)
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
		return fmt.Errorf("error creating request to create an deny list: %w", err)
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("error making a request to create an deny list: %w", err)
	}

	return nil
}

// List returns the denylist of a profile.
func (s *denylistService) List(ctx context.Context, request *ListDenylistRequest) ([]*Denylist, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), denylistAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to list the deny list: %w", err)
	}

	response := denylistResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to list the deny list: %w", err)
	}

	return response.Denylist, nil
}

// Update updates a denylist of a profile.
func (s *denylistService) Update(ctx context.Context, request *UpdateDenylistRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), denylistIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.Denylist)
	if err != nil {
		return fmt.Errorf("error creating request to update the deny list id %s: %w", request.ID, err)
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("error making a request to update the deny list id %s: %w", request.ID, err)
	}

	return nil
}

// denylistIDAPIPath returns the HTTP path for the denylist API.
func denylistIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", denylistAPIPath, id)
}
