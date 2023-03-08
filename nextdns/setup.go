package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// setupAPIPath is the HTTP path for the setup API.
const setupAPIPath = "setup"

// Setup represents the setup settings.
type Setup struct {
	Ipv4     []string       `json:"ipv4"`
	Ipv6     []string       `json:"ipv6"`
	LinkedIP *SetupLinkedIP `json:"linkedIp"`
	Dnscrypt string         `json:"dnscrypt"`
}

// GetSetupRequest encapsulates the request for getting the setup settings.
type GetSetupRequest struct {
	ProfileID string
}

// SetupService is an interface for communicating with the NextDNS setup API endpoint.
type SetupService interface {
	Get(context.Context, *GetSetupRequest) (*Setup, error)
}

// setupResponse represents the setup settings response.
type setupResponse struct {
	Setup *Setup `json:"data"`
}

// setupService represents the NextDNS setup service.
type setupService struct {
	client *Client
}

var _ SetupService = &setupService{}

// NewSetupService returns a new NextDNS setup service.
// nolint: revive
func NewSetupService(client *Client) *setupService {
	return &setupService{
		client: client,
	}
}

// Get returns the setup settings of a profile.
func (s *setupService) Get(ctx context.Context, request *GetSetupRequest) (*Setup, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), setupAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the setup settings")
	}

	response := setupResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the setup settings")
	}

	return response.Setup, nil
}
