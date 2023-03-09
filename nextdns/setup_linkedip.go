package nextdns

import (
	"context"
	"fmt"
	"net/http"
)

// setupLinkedIPAPIPath is the HTTP path for the setup linked IP API.
const setupLinkedIPAPIPath = "setup/linkedip"

type SetupLinkedIP struct {
	Servers     []string `json:"servers"`
	IP          string   `json:"ip"`
	Ddns        string   `json:"ddns"`
	UpdateToken string   `json:"updateToken"`
}

// GetSetupLinkedIPRequest encapsulates the request for getting the setup linked ip settings of a profile.
type GetSetupLinkedIPRequest struct {
	ProfileID string
}

// UpdateSetupLinkedIPRequest encapsulates the request for updating the setup linked ip of a profile.
type UpdateSetupLinkedIPRequest struct {
	ProfileID     string
	SetupLinkedIP *SetupLinkedIP
}

// SetupLinkedIPService is an interface for communicating with the NextDNS setup linked ip API endpoint.
type SetupLinkedIPService interface {
	Get(context.Context, *GetSetupLinkedIPRequest) (*SetupLinkedIP, error)
	Update(context.Context, *UpdateSetupLinkedIPRequest) error
}

// SetupLinkedIPResponse represents the setup linked ip response.
type setupLinkedIPResponse struct {
	SetupLinkedIP *SetupLinkedIP `json:"data"`
}

// SetupLinkedIPService represents the NextDNS setup linked ip service.
type setupLinkedIPService struct {
	client *Client
}

var _ SetupLinkedIPService = &setupLinkedIPService{}

// NewSetupLinkedIPService returns a new NextDNS setup linked ip service.
// nolint: revive
func NewSetupLinkedIPService(client *Client) *setupLinkedIPService {
	return &setupLinkedIPService{
		client: client,
	}
}

// Get returns the setup linked ip of a profile.
func (s *setupLinkedIPService) Get(ctx context.Context, request *GetSetupLinkedIPRequest) (*SetupLinkedIP, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), setupLinkedIPAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to get the setup linked ip settings: %w", err)
	}

	response := setupLinkedIPResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("error making a request to get the setup linked ip settings: %w", err)
	}

	return response.SetupLinkedIP, nil
}

// Update updates the setup linked ip of a profile.
func (s *setupLinkedIPService) Update(ctx context.Context, request *UpdateSetupLinkedIPRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), setupLinkedIPAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.SetupLinkedIP)
	if err != nil {
		return fmt.Errorf("error creating request to update the setup linked ip settings: %w", err)
	}

	err = s.client.do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("error making a request to update the setup linked ip settings: %w", err)
	}

	return nil
}
