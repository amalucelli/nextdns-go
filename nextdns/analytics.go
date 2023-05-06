package nextdns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	// analyticsAPIPath is the HTTP path for the denylist API.
	analyticsAPIPath = "/analytics"
	// UnidentifiedDevice is used to filter against unidentified devices.
	UnidentifiedDevice = "__UNIDENTIFIED__"
)

var ErrMissingProfile = errors.New("missing profile is required")

// analyticsService represents the NextDNS denylist service.
type analyticsService struct {
	client *Client
}

type BaseResponse struct {
	Meta struct {
		Pagination struct {
			Cursor string `json:"cursor,omitempty"`
		} `json:"pagination,omitempty"`
	} `json:"meta,omitempty"`
	Errors ErrorResponse `json:"errors,omitempty"`
}

var _ AnalyticsService = &analyticsService{}

// NewAnalyticsService returns a new NextDNS denylist service.
// nolint: revive
func NewAnalyticsService(client *Client) *analyticsService {
	return &analyticsService{
		client: client,
	}
}

type AnalyticsQuery struct {
	From   string `url:"from,omitempty"`
	To     string `url:"to,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Cursor string `url:"cursor,omitempty"`
	Device string `url:"device,omitempty"`
}

type StatusAnalyticsRequest struct {
	ProfileID string `url:"-"`
	AnalyticsQuery
}

type StatusAnalytics struct {
	Status  string `json:"status"`
	Queries int    `json:"queries"`
}

type StatusResponse struct {
	Data []*StatusAnalytics `json:"data"`

	BaseResponse
}

type DomainAnalyticsRequest struct {
	ProfileID string `url:"-"`
	Status    string `url:"status"`
	Root      bool   `url:"root"`

	AnalyticsQuery
}

type DomainsAnalytics struct {
	Domain  string `json:"domain"`
	Queries int    `json:"queries"`
	Root    string `json:"root"`
}

type DomainsAnalyticsResponse struct {
	Data []*DomainsAnalytics `json:"data"`

	BaseResponse
}

type ReasonsAnalyticsRequest struct {
	ProfileID string `url:"-"`
	AnalyticsQuery
}

type ReasonsAnalytics struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Queries int    `json:"queries"`
}

type ReasonsResponse struct {
	Data []*ReasonsAnalytics `json:"data"`

	BaseResponse
}

type IPsAnalyticsRequest struct {
	ProfileID string `url:"-"`
	AnalyticsQuery
}

type IPsAnalytics struct {
	IP      string              `json:"ip"`
	Network IPsAnalyticsNetwork `json:"network"`
	Geo     IPsAnalyticsGeo     `json:"geo"`
	Queries int                 `json:"queries"`
}

type IPsAnalyticsNetwork struct {
	Cellular bool   `json:"cellular"`
	VPN      bool   `json:"vpn"`
	ISP      string `json:"isp"`
	ASN      int    `json:"asn"`
}

type IPsAnalyticsGeo struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CountryCode string  `json:"countryCode"`
	Country     string  `json:"country"`
	City        string  `json:"city"`
}

type IPsAnalyticsResponse struct {
	Data []*IPsAnalytics `json:"data"`

	BaseResponse
}

type DevicesAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type DevicesAnalytics struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Model   string `json:"model,omitempty"`
	LocalIP string `json:"localIp,omitempty"`
	Queries int    `json:"queries"`
}

type DevicesResponse struct {
	Data []*DevicesAnalytics `json:"data"`

	BaseResponse
}

type ProtocolsAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type ProtocolsAnalytics struct {
	Protocol string `json:"protocol"`
	Queries  int    `json:"queries"`
}

type ProtocolsResponse struct {
	Data []*ProtocolsAnalytics `json:"data"`

	BaseResponse
}

type QueryTypesAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type QueryTypesAnalytics struct {
	Type    int    `url:"type,omitempty"`
	Name    string `url:"name,omitempty"`
	Queries int    `json:"queries"`
}

type QueryTypesResponse struct {
	Data []*QueryTypesAnalytics `json:"data"`

	BaseResponse
}

type IPVersionsAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type IPVersionsAnalytics struct {
	Version int `json:"version"`
	Queries int `json:"queries"`
}

type IPVersionsAnalyticsResponse struct {
	Data []*IPVersionsAnalytics `json:"data"`

	BaseResponse
}

type DNSSECAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type DNSSECAnalytics struct {
	DNSSEC  bool `json:"dnssec"`
	Queries int  `json:"queries"`
}

type DNSSECAnalyticsResponse struct {
	Data []*DNSSECAnalytics `json:"data"`

	BaseResponse
}

type EncryptionAnalyticsRequest struct {
	ProfileID string `url:"-"`

	AnalyticsQuery
}

type EncryptionAnalytics struct {
	Encrypted bool `json:"encrypted"`
	Queries   int  `json:"queries"`
}

type EncryptionResponse struct {
	Data []*EncryptionAnalytics `json:"data"`

	BaseResponse
}

type DestinationAnalyticsType string

const (
	DestinationAnalyticsTypeCountries DestinationAnalyticsType = "countries"
	DestinationAnalyticsTypeGAFAM     DestinationAnalyticsType = "gafam"
)

type DestinationsAnalyticsRequest struct {
	ProfileID string                   `url:"-"`
	Type      DestinationAnalyticsType `url:"type,omitempty"`

	AnalyticsQuery
}

type DestinationsAnalytics struct {
	Code    string   `json:"code,omitempty"`
	Domains []string `json:"domains,omitempty"`
	Company string   `json:"company,omitempty"`
	Queries int      `json:"queries"`
}

type DestinationsAnalyticsResponse struct {
	Data []*DestinationsAnalytics `json:"data"`

	BaseResponse
}

// AnalyticsService is an interface for communicating with the NextDNS denylist API endpoint.
type AnalyticsService interface {
	Status(context.Context, *StatusAnalyticsRequest) ([]*StatusAnalytics, error)
	Domains(context.Context, *DomainAnalyticsRequest) ([]*DomainsAnalytics, error)
	Reasons(context.Context, *ReasonsAnalyticsRequest) ([]*ReasonsAnalytics, error)
	IPs(context.Context, *IPsAnalyticsRequest) ([]*IPsAnalytics, error)
	Devices(context.Context, *DevicesAnalyticsRequest) ([]*DevicesAnalytics, error)
	Protocols(context.Context, *ProtocolsAnalyticsRequest) ([]*ProtocolsAnalytics, error)
	QueryTypes(context.Context, *QueryTypesAnalyticsRequest) ([]*QueryTypesAnalytics, error)
	IPVersions(context.Context, *IPVersionsAnalyticsRequest) ([]*IPVersionsAnalytics, error)
	DNSSEC(context.Context, *DNSSECAnalyticsRequest) ([]*DNSSECAnalytics, error)
	Encryption(context.Context, *EncryptionAnalyticsRequest) ([]*EncryptionAnalytics, error)
	Destinations(context.Context, *DestinationsAnalyticsRequest) ([]*DestinationsAnalytics, error)
}

func AnalyticsPath(profile, path string, query interface{}) string {
	return buildURI(profileAPIPath(profile)+analyticsAPIPath+path, query)
}

// Status returns the status analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsstatus
func (s *analyticsService) Status(ctx context.Context, query *StatusAnalyticsRequest) ([]*StatusAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}

	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/status", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var status StatusResponse
	err = s.client.do(ctx, req, &status)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return status.Data, nil
}

// Domains returns the domains analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsdomains
func (s *analyticsService) Domains(ctx context.Context, query *DomainAnalyticsRequest) ([]*DomainsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}

	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/domains", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var domains DomainsAnalyticsResponse
	err = s.client.do(ctx, req, &domains)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return domains.Data, nil
}

// Reasons returns the reasons analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsreasons
func (s *analyticsService) Reasons(ctx context.Context, query *ReasonsAnalyticsRequest) ([]*ReasonsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/reasons", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response ReasonsResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// IPs returns the IPs analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsips
func (s *analyticsService) IPs(ctx context.Context, query *IPsAnalyticsRequest) ([]*IPsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}

	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/ips", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	response := &IPsAnalyticsResponse{}
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// Devices returns the devices analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsdevices
func (s *analyticsService) Devices(ctx context.Context, query *DevicesAnalyticsRequest) ([]*DevicesAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/devices", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response DevicesResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// Protocols returns the protocols analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsprotocols
func (s *analyticsService) Protocols(ctx context.Context, query *ProtocolsAnalyticsRequest) ([]*ProtocolsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/protocols", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response ProtocolsResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// QueryTypes returns the query types analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsquerytypes
func (s *analyticsService) QueryTypes(ctx context.Context, query *QueryTypesAnalyticsRequest) ([]*QueryTypesAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/queryTypes", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response QueryTypesResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// IPVersions returns the IP versions analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsipversions
func (s *analyticsService) IPVersions(ctx context.Context, query *IPVersionsAnalyticsRequest) ([]*IPVersionsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/ipVersions", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var ipVersions IPVersionsAnalyticsResponse
	err = s.client.do(ctx, req, &ipVersions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return ipVersions.Data, nil
}

// DNSSEC returns the DNSSEC analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsdnssec
func (s *analyticsService) DNSSEC(ctx context.Context, query *DNSSECAnalyticsRequest) ([]*DNSSECAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/dnssec", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response DNSSECAnalyticsResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// Encryption returns the encryption analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsencryption
func (s *analyticsService) Encryption(ctx context.Context, query *EncryptionAnalyticsRequest) ([]*EncryptionAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/encryption", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response EncryptionResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}

// Destinations returns the destinations analytics for the given profile.
// See https://nextdns.github.io/api/#profilesprofileanalyticsdestinationstypecountries and https://nextdns.github.io/api/#profilesprofileanalyticsdestinationstypegafam
func (s *analyticsService) Destinations(ctx context.Context, query *DestinationsAnalyticsRequest) ([]*DestinationsAnalytics, error) {
	if query.ProfileID == "" {
		return nil, ErrMissingProfile
	}
	req, err := s.client.newRequest(http.MethodGet, AnalyticsPath(query.ProfileID, "/destinations", query), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMakingRequest, err)
	}

	var response DestinationsAnalyticsResponse
	err = s.client.do(ctx, req, &response)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errDoingRequest, err)
	}

	return response.Data, nil
}
