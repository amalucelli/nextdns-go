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
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
}

// GetRewritesRequest encapsulates the request for getting an rewrites.
type GetRewritesRequest struct {
	ProfileID string
}

// RewritesService is an interface for communicating with the NextDNS rewrites API endpoint.
type RewritesService interface {
	Get(context.Context, *GetRewritesRequest) ([]*Rewrites, error)
}

// rewritesResponse represents the rewrites response.
type rewritesResponse struct {
	Rewrites []*Rewrites `json:"data"`
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

// Get returns the rewrites of a profile.
func (s *rewritesService) Get(ctx context.Context, request *GetRewritesRequest) ([]*Rewrites, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), rewritesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the rewrite list")
	}

	response := rewritesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the rewrite list")
	}

	return response.Rewrites, nil
}
