package nextdns

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityTldsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/security/tlds", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
		   "data": [
			  {
				"id": "ru"
			  },
			  {
				"id": "cn"
			  },
			  {
				"id": "cf"
			  },
			  {
				"id": "accountants"
			  }
			]
		}`)
	})

	want := []*SecurityTlds{
		{ID: "ru"},
		{ID: "cn"},
		{ID: "cf"},
		{ID: "accountants"},
	}

	result, err := client.SecurityTlds.List(context.Background(), &ListSecurityTldsRequest{ProfileID: testProfileID})
	if assert.NoError(t, err, "got error when making test profile request") {
		assert.Equal(t, want, result, "got unexpected security tlds")
	}
}

func TestSecurityTldsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/profiles/%s/security/tlds", testProfileID), func(w http.ResponseWriter, r *http.Request) {
		checkHTTPMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
			"data": [
				{"id": "ru"}
			]
		}`)
	})

	err := client.SecurityTlds.Create(context.Background(), &CreateSecurityTldsRequest{ProfileID: testProfileID, SecurityTlds: []*SecurityTlds{{ID: "ru"}}})
	assert.NoError(t, err, "got error when making test security tlds request")
}
