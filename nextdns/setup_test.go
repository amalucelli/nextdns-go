package nextdns

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestSetupGet(t *testing.T) {
	c := is.New(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		out := `
{
	"data": {
		"ipv4": [
			"1.2.3.4"
		],
		"ipv6": [
			"2a07:a8c0::ab:c123",
			"2a07:a8c1::ab:c123"
		],
		"linkedIp": {
			"servers": [
				"1.1.1.1",
				"2.2.2.2"
			],
			"ip": "1.2.3.4",
			"ddns": null,
			"updateToken": "fobar"
		},
		"dnscrypt": "sdns://foobar"
	}
}`
		_, err := w.Write([]byte(out))
		c.NoErr(err)
	}))

	client, err := New(WithBaseURL(ts.URL))
	c.NoErr(err)

	ctx := context.Background()

	get, err := client.Setup.Get(ctx, &GetSetupRequest{
		ProfileID: "abc123",
	})
	want := &Setup{
		Ipv4: []string{
			"1.2.3.4",
		},
		Ipv6: []string{
			"2a07:a8c0::ab:c123",
			"2a07:a8c1::ab:c123",
		},
		LinkedIP: &SetupLinkedIP{
			Servers: []string{
				"1.1.1.1",
				"2.2.2.2",
			},
			IP:          "1.2.3.4",
			Ddns:        "",
			UpdateToken: "fobar",
		},
		Dnscrypt: "sdns://foobar",
	}

	c.NoErr(err)
	c.Equal(get, want)
}
