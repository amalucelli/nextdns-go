package nextdns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestDenylistCreate(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write([]byte(""))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()
	request := &CreateDenylistRequest{
		ProfileID: "abc123",
		Denylist: []*Denylist{
			{
				ID:     "whatsapp.net",
				Active: true,
			},
			{
				ID:     "apple.com",
				Active: false,
			},
		},
	}
	err = client.Denylist.Create(ctx, request)

	c.NoErr(err)
}

func TestDenylistGet(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		out := `{"data":[{"id":"whatsapp.net","active":true},{"id":"apple.com","active":false}]}`
		_, err := w.Write([]byte(out))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()

	list, err := client.Denylist.Get(ctx, &GetDenylistRequest{
		ProfileID: "abc123",
	})
	want := []*Denylist{
		{
			ID:     "whatsapp.net",
			Active: true,
		},
		{
			ID:     "apple.com",
			Active: false,
		},
	}

	c.NoErr(err)
	c.Equal(list, want)
}

func TestDenylistUpdate(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write([]byte(""))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()
	request := &UpdateDenylistRequest{
		ProfileID: "abc123",
		ID:        "apple.com",
		Denylist: &Denylist{
			Active: true,
		},
	}
	err = client.Denylist.Update(ctx, request)

	c.NoErr(err)
}
