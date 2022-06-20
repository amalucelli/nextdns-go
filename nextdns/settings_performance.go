package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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
	Profile string
}

// UpdateSettingsPerformanceRequest encapsulates the request for updating the settings performance of a profile.
type UpdateSettingsPerformanceRequest struct {
	Profile string
}

// SettingsPerformanceService is an interface for communicating with the NextDNS settings performance API endpoint.
type SettingsPerformanceService interface {
	Get(context.Context, *GetSettingsPerformanceRequest) (*SettingsPerformance, error)
	Update(context.Context, *UpdateSettingsPerformanceRequest, *SettingsPerformance) error
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
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsPerformanceAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the performance settings")
	}

	response := settingsPerformanceResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the performance settings")
	}

	return response.SettingsPerformance, nil
}

// Update updates the performance settings of a profile.
func (s *settingsPerformanceService) Update(ctx context.Context, request *UpdateSettingsPerformanceRequest, v *SettingsPerformance) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), settingsPerformanceAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the performance settings")
	}

	response := settingsPerformanceResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the performance settings")
	}

	return nil
}
