package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// parentalControlServicesAPIPath is the HTTP path for the parental control services API.
const parentalControlServicesAPIPath = "parentalControl/services"

// ParentalControlServices represents the parental control services of a profile.
type ParentalControlServices struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active"`
}

// CreateParentalControlServicesRequest encapsulates the request for creating a parental control services list.
type CreateParentalControlServicesRequest struct {
	Profile string
}

// UpdateParentalControlServicesRequest encapsulates the request for updating a parental control services list.
type UpdateParentalControlServicesRequest struct {
	Profile string
	ID      string
}

// GetParentalControlServicesRequest encapsulates the request for getting a parental control services list.
type GetParentalControlServicesRequest struct {
	Profile string
}

// ParentalControlServicesService is an interface for communicating with the NextDNS parental control services API endpoint.
type ParentalControlServicesService interface {
	Create(context.Context, *CreateParentalControlServicesRequest, []*ParentalControlServices) error
	Get(context.Context, *GetParentalControlServicesRequest) ([]*ParentalControlServices, error)
	Update(context.Context, *UpdateParentalControlServicesRequest, *ParentalControlServices) error
}

// parentalControlServicesResponse represents the NextDNS parental control services service.
type parentalControlServicesResponse struct {
	ParentalControlServices []*ParentalControlServices `json:"data"`
}

// parentalControlServicesService represents the NextDNS parental control services service.
type parentalControlServicesService struct {
	client *Client
}

var _ ParentalControlServicesService = &parentalControlServicesService{}

// NewParentalControlServicesService returns a new NextDNS parental control services service.
// nolint: revive
func NewParentalControlServicesService(client *Client) *parentalControlServicesService {
	return &parentalControlServicesService{
		client: client,
	}
}

// Create creates a parental control services list.
func (s *parentalControlServicesService) Create(ctx context.Context, request *CreateParentalControlServicesRequest, v []*ParentalControlServices) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), parentalControlServicesAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to create a parental control services")
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to create a parental control services")
	}

	return nil
}

// Get returns a parental control services list.
func (s *parentalControlServicesService) Get(ctx context.Context, request *GetParentalControlServicesRequest) ([]*ParentalControlServices, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), parentalControlServicesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the parental control services")
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the parental control services")
	}

	return response.ParentalControlServices, nil
}

// Update updates a parental control services list.
func (s *parentalControlServicesService) Update(ctx context.Context, request *UpdateParentalControlServicesRequest, v *ParentalControlServices) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.Profile), parentalControlServicesIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the parental control services")
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the parental control services")
	}

	return nil
}

// parentalControlServicesIDAPIPath returns the HTTP path for the parental control services API.
func parentalControlServicesIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", parentalControlServicesAPIPath, id)
}
