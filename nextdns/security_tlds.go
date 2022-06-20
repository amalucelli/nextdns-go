package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// securityTldsAPIPath is the HTTP path for the security TLDs API.
const securityTldsAPIPath = "security/tlds"

// Allowlist represents the security TLDs of a profile.
type SecurityTlds struct {
	ID string `json:"id"`
}

// CreateSecurityTldsRequest encapsulates the request for creating a security TLDs list.
type CreateSecurityTldsRequest struct {
	Profile string
}

// GetSecurityTldsRequest encapsulates the request for getting a security TLDs list.
type GetSecurityTldsRequest struct {
	Profile string
}

// SecurityTldsService is an interface for communicating with the NextDNS security TLDs API endpoint.
type SecurityTldsService interface {
	Create(context.Context, *CreateSecurityTldsRequest, []*SecurityTlds) error
	Get(context.Context, *GetSecurityTldsRequest) ([]*SecurityTlds, error)
}

// securityTldsResponse represents the security TLDs response.
type securityTldsResponse struct {
	SecurityTlds []*SecurityTlds `json:"data"`
}

// securityTldsService represents the NextDNS security TLDs service.
type securityTldsService struct {
	client *Client
}

var _ SecurityTldsService = &securityTldsService{}

// NewSecurityTldsService returns a new NextDNS security TLDs service.
// nolint: revive
func NewSecurityTldsService(client *Client) *securityTldsService {
	return &securityTldsService{
		client: client,
	}
}

// Create creates a security TLDs list.
func (s *securityTldsService) Create(ctx context.Context, request *CreateSecurityTldsRequest, v []*SecurityTlds) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), securityTldsAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to create a security tlds list")
	}

	response := securityTldsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to create a security tlds list")
	}

	return nil
}

// Get returns a security TLDs list.
func (s *securityTldsService) Get(ctx context.Context, request *GetSecurityTldsRequest) ([]*SecurityTlds, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), securityTldsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the security tlds list")
	}

	response := securityTldsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the security tlds list")
	}

	return response.SecurityTlds, nil
}
