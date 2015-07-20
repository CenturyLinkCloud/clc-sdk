package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	client := api.New(api.Config{})

	assert.NotNil(client)
}

func TestAuth(t *testing.T) {
	assert := assert.New(t)

	actualUser := &api.User{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		json.NewDecoder(r.Body).Decode(actualUser)

		fmt.Fprintf(w, `{"userName":"user@email.com","accountAlias":"ALIAS","locationAlias":"DC1","roles":["AccountAdmin","ServerAdmin"],"bearerToken":"[LONG TOKEN VALUE]"}`)
	}))
	defer ts.Close()

	config := api.Config{
		User: api.User{
			Username: "user.name",
			Password: "password",
		},
		BaseURL: ts.URL,
	}

	client := api.New(config)
	token, err := client.Auth()

	assert.Nil(err)
	assert.NotEmpty(token)
	assert.Equal(config.User.Username, actualUser.Username)
	assert.Equal(config.User.Password, actualUser.Password)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	status := "ok"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		fmt.Fprintf(w, `{"status": "%s"}`, status)
	}))
	defer ts.Close()

	client := api.New(mockConfig())
	client.Token = api.Token{Token: "valid"}

	resp := &Response{}
	err := client.Get(ts.URL, resp)

	assert.Nil(err)
	assert.Equal(status, resp.Status)
}

func TestPost(t *testing.T) {
	assert := assert.New(t)

	status := "ok"
	actualReq := &Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		json.NewDecoder(r.Body).Decode(actualReq)

		fmt.Fprintf(w, `{"status": "%s"}`, status)
	}))
	defer ts.Close()

	client := api.New(mockConfig())
	client.Token = api.Token{Token: "valid"}

	req := &Request{Status: "do stuff"}
	resp := &Response{}
	err := client.Post(ts.URL, req, resp)

	assert.Nil(err)
	assert.Equal(req, actualReq)
	assert.Equal(status, resp.Status)
}

func TestPut(t *testing.T) {
	assert := assert.New(t)

	status := "ok"
	actualReq := &Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		json.NewDecoder(r.Body).Decode(actualReq)

		fmt.Fprintf(w, `{"status": "%s"}`, status)
	}))
	defer ts.Close()

	client := api.New(mockConfig())
	client.Token = api.Token{Token: "valid"}

	req := &Request{Status: "do stuff"}
	resp := &Response{}
	err := client.Put(ts.URL, req, resp)

	assert.Nil(err)
	assert.Equal(req, actualReq)
	assert.Equal(status, resp.Status)
}

func TestPatch(t *testing.T) {
	assert := assert.New(t)

	status := "ok"
	actualReq := &Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		json.NewDecoder(r.Body).Decode(actualReq)

		fmt.Fprintf(w, `{"status": "%s"}`, status)
	}))
	defer ts.Close()

	client := api.New(mockConfig())
	client.Token = api.Token{Token: "valid"}

	req := &Request{Status: "do stuff"}
	resp := &Response{}
	err := client.Patch(ts.URL, req, resp)

	assert.Nil(err)
	assert.Equal(req, actualReq)
	assert.Equal(status, resp.Status)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	status := "ok"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			http.Error(w, "no", http.StatusMethodNotAllowed)
			return
		}

		fmt.Fprintf(w, `{"status": "%s"}`, status)
	}))
	defer ts.Close()

	client := api.New(mockConfig())
	client.Token = api.Token{Token: "valid"}

	resp := &Response{}
	err := client.Delete(ts.URL, resp)

	assert.Nil(err)
	assert.Equal(status, resp.Status)
}

func mockConfig() api.Config {
	return api.Config{
		User: api.User{
			Username: "user.name",
			Password: "password",
		},
		Alias: "t3bk",
	}
}

type Response struct {
	Status string `json:"status"`
}

type Request struct {
	Status string `json:"status"`
}
