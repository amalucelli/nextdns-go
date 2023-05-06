package nextdns

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyticsService_Status(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/status", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "status": "default",
			  "queries": 819491
			},
			{
			  "status": "blocked",
			  "queries": 132513
			},
			{
			  "status": "allowed",
			  "queries": 6923
			}
		  ]
		}`)
	})

	_, err := client.Analytics.Status(context.Background(), &StatusAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	result, err := client.Analytics.Status(context.Background(), &StatusAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test profile request") {
		assert.Equal(t, 3, len(result), "expected 3 results")
		assert.Equal(t, &StatusAnalytics{Status: "default", Queries: 819491}, result[0], "expected default status and 819491 queries")
		assert.Equal(t, &StatusAnalytics{Status: "blocked", Queries: 132513}, result[1], "expected blocked status and 132513 queries")
		assert.Equal(t, &StatusAnalytics{Status: "allowed", Queries: 6923}, result[2], "expected allowed status and 6923 queries")
	}
}

func TestAnalyticsService_Domain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/domains", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "domain": "app-measurement.com",
			  "queries": 29801
			},
			{
			  "domain": "gateway.icloud.com",
			  "root": "icloud.com",
			  "queries": 18468
			},
			{
			  "domain": "app.smartmailcloud.com",
			  "root": "smartmailcloud.com",
			  "queries": 16414
			}
		  ]
		}`)
	})

	_, err := client.Analytics.Domains(context.Background(), &DomainAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	result, err := client.Analytics.Domains(context.Background(), &DomainAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test domains request") {
		assert.Equal(t, 3, len(result), "expected 3 results")
		assert.Equal(t, &DomainsAnalytics{Domain: "app-measurement.com", Queries: 29801}, result[0], "expected app-measurement.com and 29801 queries")
		assert.Equal(t, &DomainsAnalytics{Domain: "gateway.icloud.com", Root: "icloud.com", Queries: 18468}, result[1], "expected gateway.icloud.com, icloud.com and 18468 queries")
		assert.Equal(t, &DomainsAnalytics{Domain: "app.smartmailcloud.com", Root: "smartmailcloud.com", Queries: 16414}, result[2], "expected app.smartmailcloud.com, smartmailcloud.com and 16414 queries")
	}
}

func TestAnalyticsService_Reasons(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/reasons", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "id": "blocklist:nextdns-recommended",
			  "name": "NextDNS Ads & Trackers Blocklist",
			  "queries": 131833
			},
			{
			  "id": "native:apple",
			  "name": "Native Tracking (Apple)",
			  "queries": 402
			},
			{
			  "id": "disguised-trackers",
			  "name": "Disguised Third-Party Trackers",
			  "queries": 269
			}
		  ]
		}`)
	})

	_, err := client.Analytics.Reasons(context.Background(), &ReasonsAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*ReasonsAnalytics{
		{ID: "blocklist:nextdns-recommended", Name: "NextDNS Ads & Trackers Blocklist", Queries: 131833},
		{ID: "native:apple", Name: "Native Tracking (Apple)", Queries: 402},
		{ID: "disguised-trackers", Name: "Disguised Third-Party Trackers", Queries: 269},
	}
	result, err := client.Analytics.Reasons(context.Background(), &ReasonsAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test reasons request") {
		assert.Equal(t, 3, len(result), "expected 3 results")
		assert.Equal(t, want, result, "expected reasons to match")
	}
}

func TestAnalyticsService_IPs(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/ips", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
  "data": [
    {
      "ip": "91.171.12.34",
      "network": {
        "cellular": false,
        "vpn": false,
        "isp": "Free",
        "asn": 12322
		},
      "geo": {
        "latitude": 48.8998,
        "longitude": 2.703,
        "countryCode": "FR",
        "country": "France",
        "city": "Gagny"
      },
      "queries": 136935
    },
    {
      "ip": "2a01:e0a:2cd:1234:312a:4c24:215d:185",
      "network": {
        "cellular": false,
        "vpn": false,
        "isp": "Free",
        "asn": 12322
      },
      "geo": {
        "latitude": 48.5136,
        "longitude": -1.9042,
        "countryCode": "FR",
        "country": "France",
        "city": "Miniac-Morvan"
      },
      "queries": 40410
    }
  ]
}`)
	})

	_, err := client.Analytics.IPs(context.Background(), &IPsAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*IPsAnalytics{
		{
			IP: "91.171.12.34",
			Network: IPsAnalyticsNetwork{
				Cellular: false,
				VPN:      false,
				ISP:      "Free",
				ASN:      12322,
			},
			Geo: IPsAnalyticsGeo{
				Latitude:    48.8998,
				Longitude:   2.703,
				CountryCode: "FR",
				Country:     "France",
				City:        "Gagny",
			},
			Queries: 136935,
		},
		{
			IP: "2a01:e0a:2cd:1234:312a:4c24:215d:185",
			Network: IPsAnalyticsNetwork{
				Cellular: false,
				VPN:      false,
				ISP:      "Free",
				ASN:      12322,
			},
			Geo: IPsAnalyticsGeo{
				Latitude:    48.5136,
				Longitude:   -1.9042,
				CountryCode: "FR",
				Country:     "France",
				City:        "Miniac-Morvan",
			},
			Queries: 40410,
		},
	}

	result, err := client.Analytics.IPs(context.Background(), &IPsAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test ips request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_Devices(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/devices", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
  "data": [
    {
      "id": "8TD1G",
      "name": "Romain’s iPhone",
      "model": "iPhone 12 Pro Max",
      "queries": 489885
    },
    {
      "id": "E24AR",
      "name": "MBP",
      "model": "Macbook Pro",
      "localIp": "192.168.0.11",
      "queries": 215663
    },
    {
      "id": "__UNIDENTIFIED__",
      "queries": 74242
    }
  ]
}`)
	})

	_, err := client.Analytics.Devices(context.Background(), &DevicesAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*DevicesAnalytics{
		{
			ID:      "8TD1G",
			Name:    "Romain’s iPhone",
			Model:   "iPhone 12 Pro Max",
			Queries: 489885,
		},
		{
			ID:      "E24AR",
			Name:    "MBP",
			Model:   "Macbook Pro",
			LocalIP: "192.168.0.11",
			Queries: 215663,
		},
		{
			ID:      UnidentifiedDevice,
			Queries: 74242,
		},
	}

	result, err := client.Analytics.Devices(context.Background(), &DevicesAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test devices request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_Protocols(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/protocols", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
  "data": [
    {
      "protocol": "DNS-over-HTTPS",
      "queries": 958757
    },
    {
      "protocol": "DNS-over-TLS",
      "queries": 39582
    },
    {
      "protocol": "UDP",
      "queries": 2334
    }
  ]
}`)
	})

	_, err := client.Analytics.Protocols(context.Background(), &ProtocolsAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*ProtocolsAnalytics{
		{
			Protocol: "DNS-over-HTTPS",
			Queries:  958757,
		},
		{
			Protocol: "DNS-over-TLS",
			Queries:  39582,
		},
		{
			Protocol: "UDP",
			Queries:  2334,
		},
	}

	result, err := client.Analytics.Protocols(context.Background(), &ProtocolsAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test protocols request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_QueryTypes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/queryTypes", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "type": 28,
			  "name": "AAAA",
			  "queries": 356230
			},
			{
			  "type": 1,
			  "name": "A",
			  "queries": 341812
			},
			{
			  "type": 65,
			  "name": "HTTPS",
			  "queries": 260478
			}
		  ]
		}`)
	})

	_, err := client.Analytics.QueryTypes(context.Background(), &QueryTypesAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*QueryTypesAnalytics{
		{
			Type:    28,
			Name:    "AAAA",
			Queries: 356230,
		},
		{
			Type:    1,
			Name:    "A",
			Queries: 341812,
		},
		{
			Type:    65,
			Name:    "HTTPS",
			Queries: 260478,
		},
	}

	result, err := client.Analytics.QueryTypes(context.Background(), &QueryTypesAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test query types request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_IPVersions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/ipVersions", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "version": 6,
			  "queries": 958757
			},
			{
			  "version": 4,
			  "queries": 39582
			}
		  ]
		}`)
	})

	_, err := client.Analytics.IPVersions(context.Background(), &IPVersionsAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*IPVersionsAnalytics{
		{
			Version: 6,
			Queries: 958757,
		},
		{
			Version: 4,
			Queries: 39582,
		},
	}

	result, err := client.Analytics.IPVersions(context.Background(), &IPVersionsAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test ip versions request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_DNSSEC(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/dnssec", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "queries": 958757,
			  "dnssec": true
			},
			{
			  "queries": 39582,
			  "dnssec": false
			}
		  ]
		}`)
	})

	_, err := client.Analytics.DNSSEC(context.Background(), &DNSSECAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*DNSSECAnalytics{
		{
			Queries: 958757,
			DNSSEC:  true,
		},
		{
			Queries: 39582,
			DNSSEC:  false,
		},
	}

	result, err := client.Analytics.DNSSEC(context.Background(), &DNSSECAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test dnssec request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_Encryption(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/encryption", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		  "data": [
			{
			  "queries": 958757,
			  "encrypted": true
			},
			{
			  "queries": 39582,
			  "encrypted": false
			}
		  ]
		}`)
	})

	_, err := client.Analytics.Encryption(context.Background(), &EncryptionAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	want := []*EncryptionAnalytics{
		{
			Queries:   958757,
			Encrypted: true,
		},
		{
			Queries:   39582,
			Encrypted: false,
		},
	}

	result, err := client.Analytics.Encryption(context.Background(), &EncryptionAnalyticsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test encryption request") {
		assert.Equal(t, want, result, "didn't get expected result")
	}
}

func TestAnalyticsService_Destinations(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/analytics/destinations", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		switch r.URL.Query().Get("type") {
		case "countries":
			fmt.Fprintf(w, `
				{
				  "data": [
					{
					  "code": "US",
					  "domains": [
						"app.smartmailcloud.com",
						"imap.gmail.com",
						"api.coinbase.com",
						"events-service.coinbase.com",
						"ws.coinbase.com"
					  ],
					  "queries": 209851
					},
					{
					  "code": "FR",
					  "domains": [
						"inappcheck.itunes.apple.com",
						"iphone-ld.apple.com",
						"bag.itunes.apple.com",
						"itunes.apple.com",
						"www.apple.com"
					  ],
					  "queries": 105497
					}
				  ]
				}`)
		case "gafam":
			fmt.Fprintf(w, `
				{
				  "data": [
					{
					  "company": "others",
					  "queries": 478732
					},
					{
					  "company": "apple",
					  "queries": 284832
					},
					{
					  "company": "google",
					  "queries": 159488
					}
				]
				}`)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	_, err := client.Analytics.Destinations(context.Background(), &DestinationsAnalyticsRequest{})
	assert.Equal(t, ErrMissingProfile, err, "expected missing profile error")

	wantCountries := []*DestinationsAnalytics{
		{
			Code: "US",
			Domains: []string{
				"app.smartmailcloud.com",
				"imap.gmail.com",
				"api.coinbase.com",
				"events-service.coinbase.com",
				"ws.coinbase.com",
			},
			Queries: 209851,
		},
		{
			Code: "FR",
			Domains: []string{
				"inappcheck.itunes.apple.com",
				"iphone-ld.apple.com",
				"bag.itunes.apple.com",
				"itunes.apple.com",
				"www.apple.com",
			},
			Queries: 105497,
		},
	}
	wantGafam := []*DestinationsAnalytics{
		{
			Company: "others",
			Queries: 478732,
		},
		{
			Company: "apple",
			Queries: 284832,
		},
		{
			Company: "google",
			Queries: 159488,
		},
	}

	resultCountries, err := client.Analytics.Destinations(context.Background(), &DestinationsAnalyticsRequest{ProfileID: testProfileID, Type: DestinationAnalyticsTypeCountries})
	if assert.NoError(t, err, "got error when making test destinations countries request") {
		assert.Equal(t, wantCountries, resultCountries, "didn't get expected result")
	}

	resultGafam, err := client.Analytics.Destinations(context.Background(), &DestinationsAnalyticsRequest{ProfileID: testProfileID, Type: DestinationAnalyticsTypeGAFAM})
	if assert.NoError(t, err, "got error when making test destinations gafama request") {
		assert.Equal(t, wantGafam, resultGafam, "didn't get expected result")
	}
}
