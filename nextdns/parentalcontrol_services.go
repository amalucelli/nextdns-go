package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// parentalControlServicesAPIPath is the HTTP path for the parental control services API.
const parentalControlServicesAPIPath = "parentalControl/services"

// ParentalControlServices represents the parental control services of a profile.
type ParentalControlServices struct {
	ID         string `json:"id,omitempty"`
	Active     bool   `json:"active"`
	Recreation bool   `json:"recreation"`
}

// CreateParentalControlServicesRequest encapsulates the request for creating a parental control services list.
type CreateParentalControlServicesRequest struct {
	ProfileID               string
	ParentalControlServices []*ParentalControlServices
}

// UpdateParentalControlServicesRequest encapsulates the request for updating a parental control services list.
type UpdateParentalControlServicesRequest struct {
	ProfileID               string
	ID                      string
	ParentalControlServices *ParentalControlServices
}

// GetParentalControlServicesRequest encapsulates the request for getting a parental control services list.
type ListParentalControlServicesRequest struct {
	ProfileID string
}

// ParentalControlServicesService is an interface for communicating with the NextDNS parental control services API endpoint.
type ParentalControlServicesService interface {
	Create(context.Context, *CreateParentalControlServicesRequest) error
	List(context.Context, *ListParentalControlServicesRequest) ([]*ParentalControlServices, error)
	Update(context.Context, *UpdateParentalControlServicesRequest) error
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
func (s *parentalControlServicesService) Create(ctx context.Context, request *CreateParentalControlServicesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlServicesAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.ParentalControlServices)
	if err != nil {
		return fmt.Errorf("error creating request to create a parental control services: %w", err)
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to create a parental control services: %w", err)
	}

	return nil
}

// List returns a parental control services list.
func (s *parentalControlServicesService) List(ctx context.Context, request *ListParentalControlServicesRequest) ([]*ParentalControlServices, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlServicesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to list the parental control services: %w", err)
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to list the parental control services: %w", err)
	}

	return response.ParentalControlServices, nil
}

// Update updates a parental control services list.
func (s *parentalControlServicesService) Update(ctx context.Context, request *UpdateParentalControlServicesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlServicesIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.ParentalControlServices)
	if err != nil {
		return fmt.Errorf("error creating request to update the parental control services: %w", err)
	}

	response := parentalControlServicesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to update the parental control services: %w", err)
	}

	return nil
}

// parentalControlServicesIDAPIPath returns the HTTP path for the parental control services API.
func parentalControlServicesIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", parentalControlServicesAPIPath, id)
}
