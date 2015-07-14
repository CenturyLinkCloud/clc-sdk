package aa_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/aa"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestGetAAPolicy(t *testing.T) {
	assert := assert.New(t)

	name := "12345"
	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.ID)
}

func mockStatusAPI() (*httptest.Server, *aa.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/antiAffinityPolicies/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.RequestURI, "/")
		name := parts[len(parts)-1]
		policy := aa.Response{ID: name}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policy)
	})

	mockAPI := httptest.NewServer(mux)
	config := api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: mockAPI.URL,
	}

	client := api.New(config)
	client.Token = api.Token{Token: "validtoken"}
	return mockAPI, aa.New(client)
}
