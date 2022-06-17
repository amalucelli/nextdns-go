package nextdns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const profilesAPIPath = "profiles"

type CreateProfileRequest struct {
}

type PatchProfileRequest struct {
	Profile string
}

type GetProfileRequest struct {
	Profile string
}

type ListProfileRequest struct {
}

type DeleteProfileRequest struct {
	Profile string
}

type ProfilesService interface {
	Create(context.Context, *CreateProfileRequest) (string, error)
	Get(context.Context, *GetProfileRequest) (*Profile, error)
	Patch(context.Context, *PatchProfileRequest, interface{}) error
	List(context.Context, *ListProfileRequest) ([]*Profiles, error)
	Delete(context.Context, *DeleteProfileRequest) error
}

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

type ProfileDenylist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active,omitempty"`
}

type ProfileAllowlist struct {
	ID     string `json:"id,omitempty"`
	Active bool   `json:"active,omitempty"`
}

type Profile struct {
	Name            string                 `json:"name,omitempty"`
	Security        ProfileSecurity        `json:"security,omitempty"`
	Privacy         ProfilePrivacy         `json:"privacy,omitempty"`
	ParentalControl ProfileParentalControl `json:"parentalControl,omitempty"`
	Denylist        []ProfileDenylist      `json:"denylist,omitempty"`
	Allowlist       []ProfileAllowlist     `json:"allowlist,omitempty"`
	Settings        ProfileSettings        `json:"settings,omitempty"`
}

type Profiles struct {
	ID          string `json:"id"`
	Fingerprint string `json:"fingerprint"`
	Name        string `json:"name"`
}

type profileResponse struct {
	Profile *Profile `json:"data"`
}

type profilesResponse struct {
	Profiles []*Profiles `json:"data"`
	Metadata struct {
		Pagination struct {
			Cursor string `json:"cursor"`
		} `json:"pagination"`
	} `json:"meta,omitempty"`
	Errors ErrorResponse `json:"errors,omitempty"`
}

type profilesService struct {
	client *Client
}

var _ ProfilesService = &profilesService{}

func NewProfilesService(client *Client) *profilesService {
	return &profilesService{
		client: client,
	}
}

func (ps *profilesService) List(ctx context.Context, listReq *ListProfileRequest) ([]*Profiles, error) {
	req, err := ps.client.newRequest(http.MethodGet, profilesAPIPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to list profiles")
	}

	prs := profilesResponse{}
	err = ps.client.do(ctx, req, &prs)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to list profiles")
	}

	return prs.Profiles, nil
}

func (ps *profilesService) Create(ctx context.Context, createReq *CreateProfileRequest) (string, error) {
	req, err := ps.client.newRequest(http.MethodPost, profilesAPIPath, nil)
	if err != nil {
		return "", errors.Wrap(err, "error creating request to create profile")
	}

	type newProfile struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	pr := &newProfile{}
	err = ps.client.do(ctx, req, &pr)
	if err != nil {
		return "", errors.Wrap(err, "error making a request for create profile")
	}

	return pr.Data.ID, nil
}

func (ps *profilesService) Patch(ctx context.Context, patchReq *PatchProfileRequest, v interface{}) error {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, patchReq.Profile)
	req, err := ps.client.newRequest(http.MethodPatch, path, v)
	if err != nil {
		return errors.Wrap(err, "error creating request to patch profile")
	}

	body, err := json.Marshal(&v)
	if err != nil {
		return errors.Wrap(err, "error encoding patch request")
	}

	err = ps.client.do(ctx, req, &body)
	if err != nil {
		return errors.Wrap(err, "error making a request to patch profile")
	}

	return nil
}

func (ps *profilesService) Get(ctx context.Context, getReq *GetProfileRequest) (*Profile, error) {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, getReq.Profile)
	req, err := ps.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request to get profile")
	}

	pr := profileResponse{}
	err = ps.client.do(ctx, req, &pr)
	if err != nil {
		return nil, errors.Wrap(err, "error making a request to get profile")
	}

	return pr.Profile, nil
}

func (ps *profilesService) Delete(ctx context.Context, deleteReq *DeleteProfileRequest) error {
	path := fmt.Sprintf("%s/%s", profilesAPIPath, deleteReq.Profile)
	req, err := ps.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return errors.Wrap(err, "error creating request to delete profile")
	}

	err = ps.client.do(ctx, req, nil)
	if err != nil {
		return errors.Wrap(err, "error making a request to delete profile")
	}

	return err
}
