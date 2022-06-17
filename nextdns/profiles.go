package nextdns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// profilesService is the HTTP path for the profiles API.
const profilesAPIPath = "profiles"

// CreateProfileRequest encapsulates the request for creating a new profile.
type CreateProfileRequest struct{}

// UpdateProfileRequest encapsulates the request for setting custom profile settings.
type UpdateProfileRequest struct {
	Profile string
}

// GetProfileRequest encapsulates the request for getting a profile.
type GetProfileRequest struct {
	Profile string
}

// ListProfileRequest encapsulates the request for listing all the profiles.
type ListProfileRequest struct{}

// DeleteProfileRequest encapsulates the request for deleting a profile.
type DeleteProfileRequest struct {
	Profile string
}

// ProfilesService is an interface for communicating with the NextDNS API.
type ProfilesService interface {
	Create(context.Context, *CreateProfileRequest) (string, error)
	Get(context.Context, *GetProfileRequest) (*Profile, error)
	Update(context.Context, *UpdateProfileRequest, interface{}) error
	List(context.Context, *ListProfileRequest) ([]*Profiles, error)
	Delete(context.Context, *DeleteProfileRequest) error
}

// ProfileSettings represents the settings for a profile.
type ProfileSettings struct {
	Logs struct {
		Enabled bool `json:"enabled,omitempty"`
		Drop    struct {
			IP     bool `json:"ip,omitempty"`
			Domain bool `json:"domain,omitempty"`
		} `json:"drop,omitempty"`
		Retention int    `json:"retention,omitempty"`
		Location  string `json:"location,omitempty"`
	} `json:"logs,omitempty"`
	BlockPage struct {
		Enabled bool `json:"enabled,omitempty"`
	} `json:"blockPage,omitempty"`
	Performance struct {
		Ecs             bool `json:"ecs,omitempty"`
		CacheBoost      bool `json:"cacheBoost,omitempty"`
		CnameFlattening bool `json:"cnameFlattening,omitempty"`
	} `json:"performance,omitempty"`
	Web3 bool `json:"web3,omitempty"`
}

// ProfileProfile represents the privacy settings of a provile.
type ProfilePrivacy struct {
	Blocklists []struct {
		ID string `json:"id,omitempty"`
	} `json:"blocklists,omitempty"`
	Natives []struct {
		ID string `json:"id,omitempty"`
	} `json:"natives,omitempty"`
	DisguisedTrackers bool `json:"disguisedTrackers,omitempty"`
	AllowAffiliate    bool `json:"allowAffiliate,omitempty"`
}

// ProfileSecurity represents the security settings of a profile.
type ProfileSecurity struct {
	ThreatIntelligenceFeeds bool `json:"threatIntelligenceFeeds,omitempty"`
	AiThreatDetection       bool `json:"aiThreatDetection,omitempty"`
	GoogleSafeBrowsing      bool `json:"googleSafeBrowsing,omitempty"`
	Cryptojacking           bool `json:"cryptojacking,omitempty"`
	DNSRebinding            bool `json:"dnsRebinding,omitempty"`
	IdnHomographs           bool `json:"idnHomographs,omitempty"`
	Typosquatting           bool `json:"typosquatting,omitempty"`
	Dga                     bool `json:"dga,omitempty"`
	Nrd                     bool `json:"nrd,omitempty"`
	Parking                 bool `json:"parking,omitempty"`
	Csam                    bool `json:"csam,omitempty"`
	Tlds                    []struct {
		ID string `json:"id,omitempty"`
	} `json:"tlds,omitempty"`
}

// ProfileParentalControl represents the parent control settings of a profile.
type ProfileParentalControl struct {
	Services []struct {
		ID     string `json:"id,omitempty"`
		Active bool   `json:"active,omitempty"`
	} `json:"services,omitempty"`
	Categories []struct {
		ID     string `json:"id,omitempty"`
		Active bool   `json:"active,omitempty"`
	} `json:"categories,omitempty"`
	SafeSearch            bool `json:"safeSearch,omitempty"`
	YoutubeRestrictedMode bool `json:"youtubeRestrictedMode,omitempty"`
	BlockBypass           bool `json:"blockBypass,omitempty"`
}

// ProfileDenylist represents the deny list of a profile.

type ProfileDenylist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// ProfileAllowlist represents the allow list of a profile.
type ProfileAllowlist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active,omitempty"`
}

// Profile represents a NextDNS profile.
type Profile struct {
	Name            string                 `json:"name,omitempty"`
	Security        ProfileSecurity        `json:"security,omitempty"`
	Privacy         ProfilePrivacy         `json:"privacy,omitempty"`
	ParentalControl ProfileParentalControl `json:"parentalControl,omitempty"`
	Denylist        []ProfileDenylist      `json:"denylist,omitempty"`
	Allowlist       []ProfileAllowlist     `json:"allowlist,omitempty"`
	Settings        ProfileSettings        `json:"settings,omitempty"`
}

// Profiles represents a list of NextDNS profiles.
type Profiles struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
	Name        string `json:"name"`
}

// profileResponse represents the response for the profile from the NextDNS API.
type profileResponse struct {
	Profile *Profile `json:"data"`
}

// profilesResponse represents the response for listing the profiles from the NextDNS API.
type profilesResponse struct {
	Profiles []*Profiles `json:"data"`
	Metadata struct {
		Pagination struct {
			Cursor string `json:"cursor"`
		} `json:"pagination"`
	} `json:"meta,omitempty"`
	Errors ErrorResponse `json:"errors,omitempty"`
}

// profilesService represents the NextDNS profiles service.
type profilesService struct {
	client *Client
}

var _ ProfilesService = &profilesService{}

// NewProfilesService returns a new NextDNS profiles service.
func NewProfilesService(client *Client) *profilesService {
	return &profilesService{
		client: client,
	}
}

// List returns a list of profiles.
func (ps *profilesService) List(ctx context.Context, listReq *ListProfileRequest) ([]*Profiles, error) {
	req, err := ps.client.newRequest(http.MethodGet, profilesAPIPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to list the profiles")
	}

	prs := profilesResponse{}
	err = ps.client.do(ctx, req, &prs)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to list the profiles")
	}

	return prs.Profiles, nil
}

// Create creates a profile and returns a profile ID.
func (ps *profilesService) Create(ctx context.Context, createReq *CreateProfileRequest) (string, error) {
	req, err := ps.client.newRequest(http.MethodPost, profilesAPIPath, nil)
	if err != nil {
		return "", errors.Wrap(err, "error creating request to create a profile")
	}

	type newProfile struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	pr := &newProfile{}
	err = ps.client.do(ctx, req, &pr)
	if err != nil {
		return "", errors.Wrap(err, "error making a request to create a profile")
	}

	return pr.Data.ID, nil
}

// Update updates settings of a profile.
func (ps *profilesService) Update(ctx context.Context, updateReq *UpdateProfileRequest, v interface{}) error {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, updateReq.Profile)
	req, err := ps.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to update the profile")
	}

	body, err := json.Marshal(&v)
	if err != nil {
		return errors.Wrap(err, "error encoding update request")
	}

	err = ps.client.do(ctx, req, &body)
	if err != nil {
		return errors.Wrap(err, "error making a request to update the profile")
	}

	return nil
}

// Get returns a profile.
func (ps *profilesService) Get(ctx context.Context, getReq *GetProfileRequest) (*Profile, error) {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, getReq.Profile)
	req, err := ps.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get the profile")
	}

	pr := profileResponse{}
	err = ps.client.do(ctx, req, &pr)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get the profile")
	}

	return pr.Profile, nil
}

// Delete deletes a profile.
func (ps *profilesService) Delete(ctx context.Context, deleteReq *DeleteProfileRequest) error {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, deleteReq.Profile)
	req, err := ps.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return errors.Wrap(err, "error creating request to delete the profile")
	}

	err = ps.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "error making a request to delete the profile")
	}

	return err
}
