package clc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestInitializeClient(t *testing.T) {
	assert := assert.New(t)

	client := clc.New(clc.Config{})

	assert.NotNil(client)
	assert.NotNil(client.Server)
	assert.NotNil(client.Status)
}

func mockServer(resource func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(resource))
}

func client(url string) *clc.Client {
	config := clc.Config{
		User:    clc.User{Username: "test.user", Password: "password"},
		BaseURL: url,
		Alias:   "test",
	}
	client := clc.New(config)
	client.Token = clc.Auth{Token: "validtoken"}
	return client
}
