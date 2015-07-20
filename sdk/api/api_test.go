package api_test

import (
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
