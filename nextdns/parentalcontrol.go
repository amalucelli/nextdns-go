package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// parentalControlAPIPath is the HTTP path for the parental control settings API.
const parentalControlAPIPath = "parentalControl"

// ParentalControlRecreationInterval represents the start and end time of a parental control recreation interval.
type ParentalControlRecreationInterval struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// ParentalControlRecreationTimes represents the days and times of the week when the parental control is active.
type ParentalControlRecreationTimes struct {
	Monday    *ParentalControlRecreationInterval `json:"monday,omitempty"`
	Tuesday   *ParentalControlRecreationInterval `json:"tuesday,omitempty"`
	Wednesday *ParentalControlRecreationInterval `json:"wednesday,omitempty"`
	Thursday  *ParentalControlRecreationInterval `json:"thursday,omitempty"`
	Friday    *ParentalControlRecreationInterval `json:"friday,omitempty"`
	Saturday  *ParentalControlRecreationInterval `json:"saturday,omitempty"`
	Sunday    *ParentalControlRecreationInterval `json:"sunday,omitempty"`
}

// ParentalControlRecreation represents the parental control recreation of a profile.
type ParentalControlRecreation struct {
	Times    *ParentalControlRecreationTimes `json:"times"`
	Timezone string                          `json:"timezone"`
}

// ParentalControl represents the parental control settings of a profile.
type ParentalControl struct {
	Services              []*ParentalControlServices   `json:"services,omitempty"`
	Categories            []*ParentalControlCategories `json:"categories,omitempty"`
	Recreation            *ParentalControlRecreation   `json:"recreation,omitempty"`
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
