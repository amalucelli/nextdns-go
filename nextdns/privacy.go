package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// privacyAPIPath is the HTTP path for the privacy settings API.
const privacyAPIPath = "privacy"

// Privacy represents the privacy settings of a profile.
type Privacy struct {
	Blocklists        []*PrivacyBlocklists `json:"blocklists,omitempty"`
	Natives           []*PrivacyNatives    `json:"natives,omitempty"`
	DisguisedTrackers bool                 `json:"disguisedTrackers"`
	AllowAffiliate    bool                 `json:"allowAffiliate"`
}

// UpdatePrivacyRequest encapsulates the request for updating the privacy settings of a profile.
type UpdatePrivacyRequest struct {
	ProfileID string
	Privacy   *Privacy
}

// GetPrivacyRequest encapsulates the request for getting the privacy settings of a profile.
type GetPrivacyRequest struct {
	ProfileID string
}

// PrivacyService is an interface for communicating with the NextDNS privacy settings API endpoint.
type PrivacyService interface {
	Get(context.Context, *GetPrivacyRequest) (*Privacy, error)
	Update(context.Context, *UpdatePrivacyRequest) error
}

// privacyResponse represents the NextDNS privacy settings service.
type privacyResponse struct {
	Privacy *Privacy `json:"data"`
}

// privacyService represents the NextDNS privacy settings service.
type privacyService struct {
	client *Client
}

var _ PrivacyService = &privacyService{}

// NewPrivacyService returns a new NextDNS privacy service.
// nolint: revive
func NewPrivacyService(client *Client) *privacyService {
	return &privacyService{
		client: client,
	}
}

// Get returns the privacy settings of a profile.
func (s *privacyService) Get(ctx context.Context, request *GetPrivacyRequest) (*Privacy, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get the privacy: %w", err)
	}

	response := privacyResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to get the privacy: %w", err)
	}

	return response.Privacy, nil
}

// Update updates the privacy settings of a profile.
func (s *privacyService) Update(ctx context.Context, request *UpdatePrivacyRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), privacyAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.Privacy)
	if err != nil {
		return fmt.Errorf("error creating request to update the privacy: %w", err)
	}

	response := privacyResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to update the privacy: %w", err)
	}

	return nil
}
