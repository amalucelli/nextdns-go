package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
	Profile string
}

// GetSettingsRequest encapsulates the request for getting the settings of a profile.
type GetSettingsRequest struct {
	Profile string
}

// SettingsService is an interface for communicating with the NextDNS settings API endpoint.
type SettingsService interface {
	Get(context.Context, *GetSettingsRequest) (*Settings, error)
	Update(context.Context, *UpdateSettingsRequest, *Settings) error
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
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the settings")
	}

	response := settingsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the settings")
	}

	return response.Settings, nil
}

// Update updates the settings of a profile.
func (s *settingsService) Update(ctx context.Context, request *UpdateSettingsRequest, v *Settings) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the settings")
	}

	response := settingsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the settings")
	}

	return nil
}
