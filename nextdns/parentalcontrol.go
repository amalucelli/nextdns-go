package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// parentalControlAPIPath is the HTTP path for the parental control settings API.
const parentalControlAPIPath = "parentalControl"

// ParentalControl represents the parental control settings of a profile.
type ParentalControl struct {
	Services              []*ParentalControlServices   `json:"services,omitempty"`
	Categories            []*ParentalControlCategories `json:"categories,omitempty"`
	SafeSearch            bool                         `json:"safeSearch"`
	YoutubeRestrictedMode bool                         `json:"youtubeRestrictedMode"`
	BlockBypass           bool                         `json:"blockBypass"`
}

// UpdateParentalControlRequest encapsulates the request for updating a parental control settings.
type UpdateParentalControlRequest struct {
	ProfileID       string
	ParentalControl *ParentalControl
}

// GetParentalControlRequest encapsulates the request for getting a parental control settings.
type GetParentalControlRequest struct {
	ProfileID string
}

// ParentalControlService is an interface for communicating with the NextDNS parental control API endpoint.
type ParentalControlService interface {
	Get(context.Context, *GetParentalControlRequest) (*ParentalControl, error)
	Update(context.Context, *UpdateParentalControlRequest) error
}

// parentalControlResponse represents the NextDNS parental control service.
type parentalControlResponse struct {
	ParentalControl *ParentalControl `json:"data"`
}

// parentalControlService represents the NextDNS parental control service.
type parentalControlService struct {
	client *Client
}

var _ ParentalControlService = &parentalControlService{}

// NewParentalControlService returns a new NextDNS parental control service.
// nolint: revive
func NewParentalControlService(client *Client) *parentalControlService {
	return &parentalControlService{
		client: client,
	}
}

// Get returns the parental control settings of a profile.
func (s *parentalControlService) Get(ctx context.Context, request *GetParentalControlRequest) (*ParentalControl, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the parentalControl")
	}

	response := parentalControlResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the parentalControl")
	}

	return response.ParentalControl, nil
}

// Update updates the parental control settings of a profile.
func (s *parentalControlService) Update(ctx context.Context, request *UpdateParentalControlRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.ParentalControl)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the parentalControl")
	}

	response := parentalControlResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the parentalControl")
	}

	return nil
}
