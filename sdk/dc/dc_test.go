package dc_test

import (
	"encoding/json"
	"fmt"
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

func TestGetDatacenters(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	_, err := service.GetAll()

	assert.Nil(err)
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
		if err := json.NewEncoder(w).Encode(dc); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/datacenters/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `[{"id":"dc1","name":"test datacenter","links":[{"rel":"self","href":"/v2/datacenters/test/dc1"}]}]`)
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
