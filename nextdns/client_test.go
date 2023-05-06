package nextdns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testProfileID = "testProfile"

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the API client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup(opts ...ClientOption) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	opts = append(opts, WithAPIKey("testing"))
	client, _ = New(opts...)
	client.baseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func checkHTTPMethod(t *testing.T, req *http.Request, expectedMethod string) {
	t.Helper()
	assert.Equal(t, expectedMethod, req.Method, "Expected method '%s', got %s", expectedMethod, req.Method)
}

func TestClient_Headers(t *testing.T) {
	setup()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		assert.Equal(t, "testing", r.Header.Get("X-Api-Key"))
		assert.Equal(t, contentType, r.Header.Get("Accept"))
		assert.Equal(t, userAgent, r.Header.Get("User-Agent"))
	})
	req, err := client.newRequest(http.MethodPost, "/", &UpdateSettingsRequest{})
	assert.NoError(t, err, "got error when making request")
	err = client.do(context.Background(), req, nil)
	assert.NoError(t, err, "got error when doing request")
	teardown()
}
