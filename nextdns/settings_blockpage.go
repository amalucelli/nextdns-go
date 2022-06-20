package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// settingsBlockPageAPIPath is the HTTP path for the settings block page API.
const settingsBlockPageAPIPath = "settings/blockPage"

// SettingsBlockPage represents the settings block page of a profile.
type SettingsBlockPage struct {
	Enabled bool `json:"enable"`
}

// GetSettingsBlockPageRequest encapsulates the request for getting the settings block page of a profile.
type GetSettingsBlockPageRequest struct {
	Profile string
}

// UpdateSettingsBlockPageRequest encapsulates the request for updating the settings block page of a profile.
type UpdateSettingsBlockPageRequest struct {
	Profile string
}

// SettingsBlockPageService is an interface for communicating with the NextDNS settings block page API endpoint.
type SettingsBlockPageService interface {
	Get(context.Context, *GetSettingsBlockPageRequest) (*SettingsBlockPage, error)
	Update(context.Context, *UpdateSettingsBlockPageRequest, *SettingsBlockPage) error
}

// settingsBlockPageResponse represents the settings block page response.
type settingsBlockPageResponse struct {
	SettingsBlockPage *SettingsBlockPage `json:"data"`
}

// settingsBlockPageService represents the NextDNS settings block page service.
type settingsBlockPageService struct {
	client *Client
}

var _ SettingsBlockPageService = &settingsBlockPageService{}

// NewSettingsBlockPageService returns a new NextDNS settings block page service.
// nolint: revive
func NewSettingsBlockPageService(client *Client) *settingsBlockPageService {
	return &settingsBlockPageService{
		client: client,
	}
}

// Get returns the settings block page of a profile.
func (s *settingsBlockPageService) Get(ctx context.Context, request *GetSettingsBlockPageRequest) (*SettingsBlockPage, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsBlockPageAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the block page settings")
	}

	response := settingsBlockPageResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the block page settings")
	}

	return response.SettingsBlockPage, nil
}

// Update updates the settings block page of a profile.
func (s *settingsBlockPageService) Update(ctx context.Context, request *UpdateSettingsBlockPageRequest, v *SettingsBlockPage) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsBlockPageAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the block page settings")
	}

	response := settingsBlockPageResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the block page settings")
	}

	return nil
}
