package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/server"
	"github.com/stretchr/testify/assert"
)

func TestGetServer(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	ms, service := mockServerAPI()
	defer ms.Close()

	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.Name)
}

func TestGetServerByUUID(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockServerAPI()
	defer ms.Close()

	resp, err := service.Get("5404cf5ece2042dc9f2ac16ab67416bb")

	assert.Nil(err)
	assert.Equal("va1testserver01", resp.Name)
}

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockServerAPI()
	defer ms.Close()

	server := server.Server{
		Name:           "va1testserver01",
		CPU:            1,
		MemoryGB:       1,
		GroupID:        "group",
		SourceServerID: "UBUNTU",
		Type:           "standard",
	}
	s, err := service.Create(server)

	assert.Nil(err)
	assert.True(s.IsQueued)
	assert.Equal(server.Name, s.Server)
}

func TestUpdateServer_UpdateCPU(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockServerAPI()
	defer ms.Close()

	name := "va1testserver01"
	cpu := server.CPU(1)
	resp, err := service.Update(name, cpu)

	assert.Nil(err)
	assert.Equal(name, resp.Server)
}

func TestDeleteServer(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockServerAPI()
	defer ms.Close()

	name := "va1testserver01"
	server, err := service.Delete(name)

	assert.Nil(err)
	assert.Equal(name, server.Server)
}

func mockServerAPI() (*httptest.Server, *server.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/servers/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		server := &server.Server{}
		err := json.NewDecoder(r.Body).Decode(server)
		if err != nil {
			http.Error(w, "server err", http.StatusInternalServerError)
			return
		}

		if !server.Valid() {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"server":"va1testserver01","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"},{"rel":"self","href":"/v2/servers/test/12345?uuid=True","id":"12345","verbs":["GET"]}]}`)
	})

	mux.HandleFunc("/servers/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if len(r.URL.Query()) == 0 {
				parts := strings.Split(r.RequestURI, "/")
				name := parts[len(parts)-1]
				server := &server.Response{Name: name}
				w.Header().Add("Content-Type", "application/json")
				json.NewEncoder(w).Encode(server)
				return
			}

			if r.URL.Query().Get("uuid") == "true" {
				server := &server.Response{Name: "va1testserver01"}
				w.Header().Add("Content-Type", "application/json")
				json.NewEncoder(w).Encode(server)
				return
			}
		}

		if r.Method == "DELETE" {
			parts := strings.Split(r.RequestURI, "/")
			name := parts[len(parts)-1]
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, fmt.Sprintf(`{"server":"%s","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"}]}`, name))
			return
		}

		if r.Method == "PATCH" {
			updates := make([]server.ServerUpdate, 0)
			if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			for _, v := range updates {
				if v.Op != "set" {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
			}

			parts := strings.Split(r.RequestURI, "/")
			name := parts[len(parts)-1]
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, fmt.Sprintf(`{"server":"%s","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"}]}`, name))
			return
		}

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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
	return mockAPI, server.New(client)
}
