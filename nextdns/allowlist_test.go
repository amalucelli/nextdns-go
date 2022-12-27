package nextdns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestAllowlistCreate(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write([]byte(""))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()
	request := &CreateAllowlistRequest{
		ProfileID: "abc123",
		Allowlist: []*Allowlist{
			{
				ID:     "duckduckgo.com",
				Active: true,
			},
			{
				ID:     "google.com",
				Active: false,
			},
		},
	}
	err = client.Allowlist.Create(ctx, request)

	c.NoErr(err)
}

func TestAllowlistGet(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		out := `{"data":[{"id":"google.com","active":false},{"id":"duckduckgo.com","active":true}]}`
		_, err := w.Write([]byte(out))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()

	list, err := client.Allowlist.Get(ctx, &GetAllowlistRequest{
		ProfileID: "abc123",
	})
	want := []*Allowlist{
		{
			ID:     "google.com",
			Active: false,
		},
		{
			ID:     "duckduckgo.com",
			Active: true,
		},
	}

	c.NoErr(err)
	c.Equal(list, want)
}

func TestAllowlistUpdate(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write([]byte(""))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()
	request := &UpdateAllowlistRequest{
		ProfileID: "abc123",
		ID:        "duckduckgo.com",
		Allowlist: &Allowlist{
			Active: true,
		},
	}
	err = client.Allowlist.Update(ctx, request)

	c.NoErr(err)
}
