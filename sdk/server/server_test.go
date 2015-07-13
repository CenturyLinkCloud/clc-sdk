package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/server"
	"github.com/stretchr/testify/assert"
)

func TestGetServer(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	ms := getResource(assert, name)
	defer ms.Close()

	service := service(ms.URL)
	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.Name)
}

func TestGetServerByUUID(t *testing.T) {
	assert := assert.New(t)

	name := "5404cf5ece2042dc9f2ac16ab67416bb"
	ms := getResource(assert, name)
	defer ms.Close()

	service := service(ms.URL)
	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal("va1testserver01", resp.Name)
}

func getResource(assert *assert.Assertions, name string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/servers/test/"+name, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(r.URL.Query()) == 0 {
			server := &clc.ServerResponse{Name: name}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		if r.URL.Query().Get("uuid") == "true" {
			server := &clc.ServerResponse{Name: "va1testserver01"}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
	})

	return httptest.NewServer(mux)
}

func service(url string) *server.Service {
	config := api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: url,
	}

	client := api.New(config)
	client.Token = api.Token{Token: "validtoken"}
	return server.New(client)
}
