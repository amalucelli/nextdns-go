package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// settingsAPIPath is the HTTP path for the settings API.
const settingsAPIPath = "settings"

// Settings represents the settings of a profile.
type Settings struct {
	Logs        *SettingsLogs        `json:"logs,omitempty"`
	BlockPage   *SettingsBlockPage   `json:"blockPage,omitempty"`
	Performance *SettingsPerformance `json:"performance,omitempty"`
	Web3        bool                 `json:"web3"`
}

// UpdateSettingsRequest encapsulates the request for updating the settings of a profile.
type UpdateSettingsRequest struct {
	ProfileID string
	Settings  *Settings
}

// GetSettingsRequest encapsulates the request for getting the settings of a profile.
type GetSettingsRequest struct {
	ProfileID string
}

// SettingsService is an interface for communicating with the NextDNS settings API endpoint.
type SettingsService interface {
	Get(context.Context, *GetSettingsRequest) (*Settings, error)
	Update(context.Context, *UpdateSettingsRequest) error
}

// settingsResponse represents the settings response.
type settingsResponse struct {
	Settings *Settings `json:"data"`
}

// settingsService represents the NextDNS settings service.
type settingsService struct {
	client *Client
}

var _ SettingsService = &settingsService{}

// NewSettingsService returns a new NextDNS settings service.
// nolint: revive
func NewSettingsService(client *Client) *settingsService {
	return &settingsService{
		client: client,
	}
}

// Get returns the settings of a profile.
func (s *settingsService) Get(ctx context.Context, request *GetSettingsRequest) (*Settings, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get the settings: %w", err)
	}

	response := settingsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to get the settings: %w", err)
	}

	return response.Settings, nil
}

// Update updates the settings of a profile.
func (s *settingsService) Update(ctx context.Context, request *UpdateSettingsRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.Settings)
	if err != nil {
		return fmt.Errorf("error creating request to update the settings: %w", err)
	}

	response := settingsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to update the settings: %w", err)
	}

	return nil
}
