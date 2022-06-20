package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// parentalControlCategoriesAPIPath is the HTTP path for the parental control categories API.
const parentalControlCategoriesAPIPath = "parentalControl/categories"

// ParentalControlCategories represents the parental control categories of a profile.
type ParentalControlCategories struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active"`
}

// CreateParentalControlCategoriesRequest encapsulates the request for creating a parental control categories list.
type CreateParentalControlCategoriesRequest struct {
	ProfileID                 string
	ParentalControlCategories []*ParentalControlCategories
}

// UpdateParentalControlCategoriesRequest encapsulates the request for updating a parental control categories list.
type UpdateParentalControlCategoriesRequest struct {
	ProfileID                 string
	ID                        string
	ParentalControlCategories *ParentalControlCategories
}

// GetParentalControlCategoriesRequest encapsulates the request for getting a parental control categories list.
type GetParentalControlCategoriesRequest struct {
	ProfileID string
}

// ParentalControlCategoriesService is an interface for communicating with the NextDNS parental control categories API endpoint.
type ParentalControlCategoriesService interface {
	Create(context.Context, *CreateParentalControlCategoriesRequest) error
	Get(context.Context, *GetParentalControlCategoriesRequest) ([]*ParentalControlCategories, error)
	Update(context.Context, *UpdateParentalControlCategoriesRequest) error
}

// parentalControlCategoriesResponse represents the parental control categories response.
type parentalControlCategoriesResponse struct {
	ParentalControlCategories []*ParentalControlCategories `json:"data"`
}

// parentalControlCategoriesService represents the NextDNS parental control categories service.
type parentalControlCategoriesService struct {
	client *Client
}

var _ ParentalControlCategoriesService = &parentalControlCategoriesService{}

// NewParentalControlCategoriesService returns a new NextDNS parental control categories service.
// nolint: revive
func NewParentalControlCategoriesService(client *Client) *parentalControlCategoriesService {
	return &parentalControlCategoriesService{
		client: client,
	}
}

// Create creates a parental control categories list.
func (s *parentalControlCategoriesService) Create(ctx context.Context, request *CreateParentalControlCategoriesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlCategoriesAPIPath)
	req, err := s.client.newRequest(http.MethodPut, path, request.ParentalControlCategories)
	if err != nil {
		return errors.Wrap(err, "error creating request to create a parental control categories")
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to create a parental control categories")
	}

	return nil
}

// Get returns a parental control categories list.
func (s *parentalControlCategoriesService) Get(ctx context.Context, request *GetParentalControlCategoriesRequest) ([]*ParentalControlCategories, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlCategoriesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the parental control categories")
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the parental control categories")
	}

	return response.ParentalControlCategories, nil
}

// Update updates a parental control categories list.
func (s *parentalControlCategoriesService) Update(ctx context.Context, request *UpdateParentalControlCategoriesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlCategoriesIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.ParentalControlCategories)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the parental control categories")
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the parental control categories")
	}

	return nil
}

// parentalControlCategoriesIDAPIPath returns the HTTP path for the parental control categories API.
func parentalControlCategoriesIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", parentalControlCategoriesAPIPath, id)
}
