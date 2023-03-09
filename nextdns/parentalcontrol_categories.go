package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// parentalControlCategoriesAPIPath is the HTTP path for the parental control categories API.
const parentalControlCategoriesAPIPath = "parentalControl/categories"

// ParentalControlCategories represents the parental control categories of a profile.
type ParentalControlCategories struct {
	ID         string `json:"id,omitempty"`
	Active     bool   `json:"active"`
	Recreation bool   `json:"recreation"`
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

// ListParentalControlCategoriesRequest encapsulates the request for getting a parental control categories list.
type ListParentalControlCategoriesRequest struct {
	ProfileID string
}

// ParentalControlCategoriesService is an interface for communicating with the NextDNS parental control categories API endpoint.
type ParentalControlCategoriesService interface {
	Create(context.Context, *CreateParentalControlCategoriesRequest) error
	List(context.Context, *ListParentalControlCategoriesRequest) ([]*ParentalControlCategories, error)
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
		return fmt.Errorf("error creating request to create a parental control categories: %w", err)
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to create a parental control categories: %w", err)
	}

	return nil
}

// List returns a parental control categories list.
func (s *parentalControlCategoriesService) List(ctx context.Context, request *ListParentalControlCategoriesRequest) ([]*ParentalControlCategories, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlCategoriesAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to list the parental control categories: %w", err)
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to list the parental control categories: %w", err)
	}

	return response.ParentalControlCategories, nil
}

// Update updates a parental control categories list.
func (s *parentalControlCategoriesService) Update(ctx context.Context, request *UpdateParentalControlCategoriesRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), parentalControlCategoriesIDAPIPath(request.ID))
	req, err := s.client.newRequest(http.MethodPatch, path, request.ParentalControlCategories)
	if err != nil {
		return fmt.Errorf("error creating request to update the parental control categories: %w", err)
	}

	response := parentalControlCategoriesResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return fmt.Errorf("error making a request to update the parental control categories: %w", err)
	}

	return nil
}

// parentalControlCategoriesIDAPIPath returns the HTTP path for the parental control categories API.
func parentalControlCategoriesIDAPIPath(id string) string {
	return fmt.Sprintf("%s/%s", parentalControlCategoriesAPIPath, id)
}
