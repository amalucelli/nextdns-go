package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// settingsLogsAPIPath is the HTTP path for the settings logs API.
const settingsLogsAPIPath = "settings/logs"

// SettingsLogsDrop represents the settings logs privacy adjustments of a profile.
type SettingsLogsDrop struct {
	IP     bool `json:"ip"`
	Domain bool `json:"domain"`
}

// SettingsLogs represents the settings logs of a profile.
type SettingsLogs struct {
	Enabled   bool              `json:"enabled"`
	Drop      *SettingsLogsDrop `json:"drop,omitempty"`
	Retention int               `json:"retention,omitempty"`
	Location  string            `json:"location,omitempty"`
}

// GetSettingsLogsRequest encapsulates the request for getting the settings logs of a profile.
type GetSettingsLogsRequest struct {
	ProfileID string
}

// UpdateSettingsLogsRequest encapsulates the request for updating the settings logs of a profile.
type UpdateSettingsLogsRequest struct {
	ProfileID    string
	SettingsLogs *SettingsLogs
}

// SettingsLogsService is an interface for communicating with the NextDNS settings logs API endpoint.
type SettingsLogsService interface {
	Get(context.Context, *GetSettingsLogsRequest) (*SettingsLogs, error)
	Update(context.Context, *UpdateSettingsLogsRequest) error
}

// settingsLogsResponse represents the settings logs response.
type settingsLogsResponse struct {
	SettingsLogs *SettingsLogs `json:"data"`
}

// settingsLogsService represents the NextDNS settings logs service.
type settingsLogsService struct {
	client *Client
}

var _ SettingsLogsService = &settingsLogsService{}

// NewSettingsLogsService returns a new NextDNS settings logs service.
// nolint: revive
func NewSettingsLogsService(client *Client) *settingsLogsService {
	return &settingsLogsService{
		client: client,
	}
}

// Get returns the settings logs of a profile.
func (s *settingsLogsService) Get(ctx context.Context, request *GetSettingsLogsRequest) (*SettingsLogs, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsLogsAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the logs settings")
	}

	response := settingsLogsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the logs settings")
	}

	return response.SettingsLogs, nil
}

// Update updates the settings logs of a profile.
func (s *settingsLogsService) Update(ctx context.Context, request *UpdateSettingsLogsRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsLogsAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.SettingsLogs)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the logs settings")
	}

	response := settingsLogsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the logs settings")
	}

	return nil
}
