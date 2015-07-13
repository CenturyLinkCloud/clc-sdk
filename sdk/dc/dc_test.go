package dc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/dc"
	"github.com/stretchr/testify/assert"
)

func TestGetDatacenter(t *testing.T) {
	assert := assert.New(t)

	name := "va1"
	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.ID)
}

func mockStatusAPI() (*httptest.Server, *dc.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/datacenters/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.RawQuery != "groupLinks=true" {
			http.Error(w, "bad reqest", http.StatusBadRequest)
			return
		}

		parts := strings.Split(r.RequestURI, "/")
		name := strings.Split(parts[len(parts)-1], "?")[0]
		dc := dc.Response{ID: name}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dc)
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
	return mockAPI, dc.New(client)
}
