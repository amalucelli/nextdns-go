package nextdns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// securityAPIPath is the HTTP path for the security API.
const securityAPIPath = "security"

// Security represents the security settings of a profile.
type Security struct {
	ThreatIntelligenceFeeds bool            `json:"threatIntelligenceFeeds"`
	AiThreatDetection       bool            `json:"aiThreatDetection"`
	GoogleSafeBrowsing      bool            `json:"googleSafeBrowsing"`
	Cryptojacking           bool            `json:"cryptojacking"`
	DNSRebinding            bool            `json:"dnsRebinding"`
	IdnHomographs           bool            `json:"idnHomographs"`
	Typosquatting           bool            `json:"typosquatting"`
	Dga                     bool            `json:"dga"`
	Nrd                     bool            `json:"nrd"`
	DDNS                    bool            `json:"ddns"`
	Parking                 bool            `json:"parking"`
	Csam                    bool            `json:"csam"`
	Tlds                    []*SecurityTlds `json:"tlds,omitempty"`
}

// UpdateSecurityRequest encapsulates the request for updating security settings.
type UpdateSecurityRequest struct {
	ProfileID string
	Security  *Security
}

// GetSecurityRequest encapsulates the request for getting a security settings.
type GetSecurityRequest struct {
	ProfileID string
}

// SecurityService is an interface for communicating with the NextDNS security API endpoint.
type SecurityService interface {
	Get(context.Context, *GetSecurityRequest) (*Security, error)
	Update(context.Context, *UpdateSecurityRequest) error
}

// securityResponse represents the security settings response.
type securityResponse struct {
	Security *Security `json:"data"`
}

// securityService represents the NextDNS security service.
type securityService struct {
	client *Client
}

var _ SecurityService = &securityService{}

// NewSecurityService returns a new NextDNS security service.
// nolint: revive
func NewSecurityService(client *Client) *securityService {
	return &securityService{
		client: client,
	}
}

// Get returns the security settings of a profile.
func (s *securityService) Get(ctx context.Context, request *GetSecurityRequest) (*Security, error) {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), securityAPIPath)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the security settings")
	}

	response := securityResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the security settings")
	}

	return response.Security, nil
}

// Update updates the security settings of a profile.
func (s *securityService) Update(ctx context.Context, request *UpdateSecurityRequest) error {
	path := fmt.Sprintf("%s/%s", profileAPIPath(request.ProfileID), securityAPIPath)
	req, err := s.client.newRequest(http.MethodPatch, path, request.Security)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the security settings")
	}

	response := securityResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the security settings")
	}

	return nil
}
