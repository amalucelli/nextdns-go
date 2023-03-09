package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// settingsPerformanceAPIPath is the HTTP path for the settings performance API.
const settingsPerformanceAPIPath = "settings/performance"

// SettingsPerformance represents the settings performance of a profile.
type SettingsPerformance struct {
	Ecs             bool `json:"ecs"`
	CacheBoost      bool `json:"cacheBoost"`
	CnameFlattening bool `json:"cnameFlattening"`
}

// GetSettingsPerformanceRequest encapsulates the request for getting the settings performance of a profile.
type GetSettingsPerformanceRequest struct {
	ProfileID string
}

// UpdateSettingsPerformanceRequest encapsulates the request for updating the settings performance of a profile.
type UpdateSettingsPerformanceRequest struct {
	ProfileID           string
	SettingsPerformance *SettingsPerformance
}

// SettingsPerformanceService is an interface for communicating with the NextDNS settings performance API endpoint.
type SettingsPerformanceService interface {
	Get(context.Context, *GetSettingsPerformanceRequest) (*SettingsPerformance, error)
	Update(context.Context, *UpdateSettingsPerformanceRequest) error
}

// settingsPerformanceResponse represents the settings performance response.
type settingsPerformanceResponse struct {
	SettingsPerformance *SettingsPerformance `json:"data"`
}

// settingsPerformanceService represents the NextDNS settings performance service.
type settingsPerformanceService struct {
	client *Client
}

var _ SettingsPerformanceService = &settingsPerformanceService{}

// NewSettingsPerformanceService returns a new NextDNS settings performance service.
// nolint: revive
func NewSettingsPerformanceService(client *Client) *settingsPerformanceService {
	return &settingsPerformanceService{
		client: client,
	}
}

// Get returns the performance settings of a profile.
func (s *settingsPerformanceService) Get(ctx context.Context, request *GetSettingsPerformanceRequest) (*SettingsPerformance, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsPerformanceAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get the performance settings: %w", err)
	}

	response := settingsPerformanceResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to get the performance settings: %w", err)
	}

	return response.SettingsPerformance, nil
}

// Update updates the performance settings of a profile.
func (s *settingsPerformanceService) Update(ctx context.Context, request *UpdateSettingsPerformanceRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), settingsPerformanceAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.SettingsPerformance)
	if err != nil {
		return fmt.Errorf("error creating request to update the performance settings: %w", err)
	}

	response := settingsPerformanceResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to update the performance settings: %w", err)
	}

	return nil
}
