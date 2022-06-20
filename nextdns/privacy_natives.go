package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// privacyNativesAPIPath is the HTTP path for the privacy native tracking protection API.
const privacyNativesAPIPath = "privacy/natives"

// PrivacyNatives represents a privacy native tracking protection of a profile.
type PrivacyNatives struct {
	ID string `json:"id"`
}

// CreatePrivacyNativesRequest encapsulates the request for creating a privacy native tracking protection list.
type CreatePrivacyNativesRequest struct {
	ProfileID      string
	PrivacyNatives []*PrivacyNatives
}

// GetPrivacyNativesRequest encapsulates the request for getting the privacy native tracking protection list.
type GetPrivacyNativesRequest struct {
	ProfileID string
}

// PrivacyNativesService is an interface for communicating with the NextDNS privacy native tracking protection API endpoint.
type PrivacyNativesService interface {
	Create(context.Context, *CreatePrivacyNativesRequest) error
	Get(context.Context, *GetPrivacyNativesRequest) ([]*PrivacyNatives, error)
}

// privacyNativesResponse represents the NextDNS privacy native tracking protection service.
type privacyNativesResponse struct {
	PrivacyNatives []*PrivacyNatives `json:"data"`
}

// privacyNativesService represents the NextDNS privacy native tracking protection service.
type privacyNativesService struct {
	client *Client
}

var _ PrivacyNativesService = &privacyNativesService{}

// NewPrivacyNativesService returns a new NextDNS privacy native tracking protection service.
// nolint: revive
func NewPrivacyNativesService(client *Client) *privacyNativesService {
	return &privacyNativesService{
		client: client,
	}
}

// Create creates a privacy native tracking protection list.
func (s *privacyNativesService) Create(ctx context.Context, request *CreatePrivacyNativesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyNativesAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.PrivacyNatives)
	if err != nil {
		return errors.Wrap(err, "error creating request to create a privacy native list")
	}

	response := privacyNativesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to create a privacy native list")
	}

	return nil
}

// Get returns the privacy native tracking protection list.
func (s *privacyNativesService) Get(ctx context.Context, request *GetPrivacyNativesRequest) ([]*PrivacyNatives, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyNativesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the privacy native list")
	}

	response := privacyNativesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the privacy native list")
	}

	return response.PrivacyNatives, nil
}
