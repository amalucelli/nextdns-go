package nextdns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestSetupLinkedIpGet(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		out := `
{
	"data": {
		"servers": [
			"1.1.1.1",
			"2.2.2.2"
		],
		"ip": "1.2.3.4",
		"ddns": null,
		"updateToken": "fobar"
	}
}`
		_, err := w.Write([]byte(out))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()

	get, err := client.SetupLinkedIP.Get(ctx, &GetSetupLinkedIPRequest{
		ProfileID: "abc123",
	})
	want := &SetupLinkedIP{
		Servers: []string{
			"1.1.1.1",
			"2.2.2.2",
		},
		IP:          "1.2.3.4",
		Ddns:        "",
		UpdateToken: "fobar",
	}

	c.NoErr(err)
	c.Equal(get, want)
}

func TestSetupLinkedIpUpdate(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		_, err := w.Write([]byte(""))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL), WithDebug())
	c.NoErr(err)

	ctx := context.Background()
	request := &UpdateSetupLinkedIPRequest{
		ProfileID: "abc123",
		SetupLinkedIP: &SetupLinkedIP{
			Ddns: "foobar.no-ip.org",
		},
	}
	err = client.SetupLinkedIP.Update(ctx, request)

	c.NoErr(err)
}
