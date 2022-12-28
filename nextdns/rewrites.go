package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// rewritesAPIPath is the HTTP path for the rewrites API.
const rewritesAPIPath = "rewrites"

// Rewrites represents the rewrite list of a profile.
type Rewrites struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content"`
}

// CreateRewritesRequest encapsulates the request for creating a new rewrite.
type CreateRewritesRequest struct {
	ProfileID string
	Rewrites  *Rewrites
}

// ListRewritesRequest encapsulates the request for getting an rewrites.
type ListRewritesRequest struct {
	ProfileID string
}

// DeleteRewritesRequest encapsulates the request for deleting a rewrite.
type DeleteRewritesRequest struct {
	ProfileID string
	ID        string
}

// RewritesService is an interface for communicating with the NextDNS rewrites API endpoint.
type RewritesService interface {
	Create(context.Context, *CreateRewritesRequest) (string, error)
	List(context.Context, *ListRewritesRequest) ([]*Rewrites, error)
	Delete(context.Context, *DeleteRewritesRequest) error
}

// rewritesResponse represents the rewrites response.
type rewritesResponse struct {
	Rewrites []*Rewrites `json:"data"`
}

// createRewritesResponse represents the response when creating a rewrite from the NextDNS API.
type createRewritesResponse struct {
	Rewrites *Rewrites `json:"data"`
}

// privacyService represents the NextDNS rewrites service.
type rewritesService struct {
	client *Client
}

var _ RewritesService = &rewritesService{}

// NewRewritesService returns a new NextDNS rewrites service.
// nolint: revive
func NewRewritesService(client *Client) *rewritesService {
	return &rewritesService{
		client: client,
	}
}

// Create creates a rewrite and returns its ID.
func (s *rewritesService) Create(ctx context.Context, request *CreateRewritesRequest) (string, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), rewritesAPIPath)

	req, err := s.client.newRequest(http.MethodPost, path, request.Rewrites)
	if err != nil {
		return "", errors.Wrap(err, "error creating request to create a rewrite")
	}

	response := &createRewritesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return "", errors.Wrap(err, "error making a request to create a rewrite")
	}

	return response.Rewrites.ID, nil
}

// List returns the rewrites of a profile.
func (s *rewritesService) List(ctx context.Context, request *ListRewritesRequest) ([]*Rewrites, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), rewritesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to list the rewrite list")
	}

	response := rewritesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to list the rewrite list")
	}

	return response.Rewrites, nil
}

// Delete deletes a profile.
func (s *rewritesService) Delete(ctx context.Context, request *DeleteRewritesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), rewritesIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return errors.Wrap(err, "error creating request to delete the rewrite")
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "error making a request to delete the rewrite")
	}

	return err
}

// rewritesIDAPIPath returns the HTTP path for the rewrites API.
func rewritesIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", rewritesAPIPath, id)
}
