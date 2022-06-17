package nextdns

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-cleanhttp"
)

const (
	baseURL     = "https://api.nextdns.io/"
	contentType = "application/json"
)

type Client struct {
	client  *http.Client
	baseURL *url.URL

	Profiles ProfilesService
}

type ClientOption func(c *Client) error

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

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		if apiKey == "" {
			return errors.New(errEmptyAPIToken)
		}

		transport := authTransport{
			rt:     c.client.Transport,
			apiKey: apiKey,
		}

		c.client.Transport = &transport
		return nil
	}
}

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		if client == nil {
			client = cleanhttp.DefaultClient()
		}

		c.client = client
		return nil
	}
}

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

	c.Profiles = &profilesService{client: c}

	return c, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return c.handleResponse(ctx, res, v)
}

func (c *Client) handleResponse(ctx context.Context, res *http.Response, v interface{}) error {
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	meta := map[string]string{
		"body":        string(out),
		"http_status": http.StatusText(res.StatusCode),
	}

	if res.StatusCode >= http.StatusBadRequest {
		if res.StatusCode >= http.StatusInternalServerError {
			return &Error{
				Type:    ErrorTypeServiceError,
				Message: errInternalServiceError,
				Errors:  nil,
				Meta:    meta,
			}
		}

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

		var errType ErrorType

		switch res.StatusCode {
		case http.StatusForbidden:
			errType = ErrorTypeAuthentication
		case http.StatusNotFound:
			errType = ErrorTypeNotFound
		default:
			errType = ErrorTypeRequest
		}

		return &Error{
			Type:    errType,
			Message: errResponseError,
			Errors:  errorRes,
			Meta:    meta,
		}
	}

	if v == nil {
		return nil
	}

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

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet:
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

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", contentType)
	return req, nil
}

type authTransport struct {
	rt     http.RoundTripper
	apiKey string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Api-Key", t.apiKey)
	return t.rt.RoundTrip(req)
}
