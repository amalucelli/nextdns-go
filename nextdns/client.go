package nextdns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	baseURL     = "https://api.nextdns.io/"
	contentType = "application/json"
	userAgent   = "nextdns-go"
)

// Client represents a NextDNS client.
type Client struct {
	client  *http.Client
	baseURL *url.URL

	// Service for the Profile.
	Profiles ProfilesService

	// Services for the Allowlist and Denylist.
	Allowlist AllowlistService
	Denylist  DenylistService

	// Services for the ParentalControl.
	ParentalControlServices   ParentalControlServicesService
	ParentalControlCategories ParentalControlCategoriesService

	// Services for the Privacy.
	Privacy           PrivacyService
	PrivacyBlocklists PrivacyBlocklistsService
	PrivacyNatives    PrivacyNativesService

	// Services for the Settings.
	Settings            SettingsService
	SettingsLogs        SettingsLogsService
	SettingsBlockPage   SettingsBlockPageService
	SettingsPerformance SettingsPerformanceService

	// Services for the Security.
	Security     SecurityService
	SecurityTlds SecurityTldsService

	// Services for the Rewrites.
	Rewrites RewritesService

	// Debug mode for the HTTP requests.
	Debug bool
}

// ClientOption is a function that can be used to customize the client.
type ClientOption func(c *Client) error

// WithBaseURL sets the base URL of the NextDNS API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		c.baseURL = parsedURL
		return nil
	}
}

// WithAPIKey sets the API key to be used for requests.
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		if apiKey == "" {
			return ErrEmptyAPIToken
		}

		transport := authTransport{
			rt:     c.client.Transport,
			apiKey: apiKey,
		}

		c.client.Transport = &transport
		return nil
	}
}

// WithDebug enables debug mode.
func WithDebug() ClientOption {
	return func(c *Client) error {
		c.Debug = true
		return nil
	}
}

// WithHTTPClient sets a custom HTTP client that can be used for requests.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		if client == nil {
			client = cleanhttp.DefaultClient()
		}

		c.client = client
		return nil
	}
}

// New instantiates a new NextDNS client.
func New(opts ...ClientOption) (*Client, error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:  cleanhttp.DefaultClient(),
		baseURL: baseURL,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// Initialize the services for the Profile.
	c.Profiles = NewProfilesService(c)

	// Initialize the services for the Allowlist and Denylist.
	c.Allowlist = NewAllowlistService(c)
	c.Denylist = NewDenylistService(c)

	// Initialize the services for the ParentalControl.
	c.ParentalControlServices = NewParentalControlServicesService(c)
	c.ParentalControlCategories = NewParentalControlCategoriesService(c)

	// Initialize the services for the Privacy.
	c.Privacy = NewPrivacyService(c)
	c.PrivacyBlocklists = NewPrivacyBlocklistsService(c)
	c.PrivacyNatives = NewPrivacyNativesService(c)

	// Initialize the services for the Settings.
	c.Settings = NewSettingsService(c)
	c.SettingsLogs = NewSettingsLogsService(c)
	c.SettingsBlockPage = NewSettingsBlockPageService(c)
	c.SettingsPerformance = NewSettingsPerformanceService(c)

	// Initialize the services for the Security.
	c.Security = NewSecurityService(c)
	c.SecurityTlds = NewSecurityTldsService(c)

	// Initialize the services for the Rewrites.
	c.Rewrites = NewRewritesService(c)

	return c, nil
}

// do executes an HTTP request and decodes the response into v.
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return c.handleResponse(ctx, res, v)
}

// handleResponse handles the response from the NextDNS API and decodes the response into v if provided.
// The goal is to handle the common errors that can occur when making a request to the NextDNS API,
// and also provide custom error responses for the client.
func (c *Client) handleResponse(ctx context.Context, res *http.Response, v interface{}) error {
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if c.Debug {
		if string(out) == "" {
			fmt.Printf("[DEBUG] RESPONSE: StatusCode:%d\n", res.StatusCode)
		} else {
			fmt.Printf("[DEBUG] RESPONSE: StatusCode:%d, Body:%v\n", res.StatusCode, string(out))
		}
	}

	// If there is no response body, then we don't need to do anything.
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	// Sets some default additional informations that can be used by the client to debug the error.
	meta := map[string]string{
		"body":        string(out),
		"http_status": http.StatusText(res.StatusCode),
	}

	// If the response is not a 200, then we need to handle the error.
	// TODO(amalucelli): Report the behavior to NextDNS, but there are errors that return HTTP 200 ("duplicate" case).
	if res.StatusCode >= http.StatusBadRequest || strings.Contains(string(out), "\"errors\"") {
		if res.StatusCode >= http.StatusInternalServerError {
			return &Error{
				Type:    ErrorTypeServiceError,
				Message: errInternalServiceError,
				Errors:  nil,
				Meta:    meta,
			}
		}

		// Tries to handle the error response body from the NextDNS API,
		// encapsulated in a client error.
		errorRes := &ErrorResponse{}
		err = json.Unmarshal(out, errorRes)
		if err != nil {
			var jsonErr *json.SyntaxError
			if errors.As(err, &jsonErr) {
				meta["err"] = jsonErr.Error()
				return &Error{
					Type:    ErrorTypeMalformed,
					Message: errMalformedErrorBody,
					Errors:  nil,
					Meta:    meta,
				}
			}
			return err
		}

		// Sets custom error messages for the client based on the HTTP status code.
		var errType ErrorType

		switch res.StatusCode {
		case http.StatusForbidden:
			errType = ErrorTypeAuthentication
		case http.StatusNotFound:
			errType = ErrorTypeNotFound
		default:
			errType = ErrorTypeRequest
		}

		// Returns the error response from the NextDNS API encapsulated in a client error.
		return &Error{
			Type:    errType,
			Message: errResponseError,
			Errors:  errorRes,
			Meta:    meta,
		}
	}

	// Returns if there is no object to decode.
	if v == nil {
		return nil
	}

	// Decodes the response body into the provided object.
	err = json.Unmarshal(out, &v)
	if err != nil {
		var jsonErr *json.SyntaxError
		if errors.As(err, &jsonErr) {
			meta["err"] = jsonErr.Error()
			return &Error{
				Type:    ErrorTypeMalformed,
				Message: errMalformedError,
				Errors:  nil,
				Meta:    meta,
			}
		}
		return err
	}

	return nil
}

// newRequest creates a new HTTP request.
func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet:
		if c.Debug {
			fmt.Printf("[DEBUG] REQUEST: Method:%s, URL:%s\n", method, u.String())
		}
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		buf := new(bytes.Buffer)
		if body != nil {
			err = json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		}
		if c.Debug {
			fmt.Printf("[DEBUG] REQUEST: Method:%s, URL:%s, Body:%v", method, u.String(), buf.String())
		}
		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}

// authHeader represents a RoundTripper that adds an authorization header to the request.
type authTransport struct {
	rt     http.RoundTripper
	apiKey string
}

// RoundTrip adds the authorization header to requests.
func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Api-Key", t.apiKey)
	return t.rt.RoundTrip(req)
}
