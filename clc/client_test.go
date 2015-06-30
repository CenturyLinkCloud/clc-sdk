package clc_test

import (
	"encoding/json"
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
}

func TestAuth(t *testing.T) {
	assert := assert.New(t)

	expectedToken := "TOKEN"
	resource := authResource(assert, expectedToken)
	ms := mockServer(resource)
	defer ms.Close()

	client := client(ms.URL)
	token, err := client.Auth()

	assert.Nil(err)
	assert.Equal(expectedToken, token)
}

func authResource(assert *assert.Assertions, token string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			assert.Fail("Authentication method should be POST", r.Method)
		}

		if r.URL.Path != "/authentication/login" {
			assert.Fail("Authentication hitting wrong endpoint", r.URL.Path)
		}

		user := &clc.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			assert.Fail("Decoding user failed.")
		}

		auth := &clc.Auth{
			Username: user.Username,
			Alias:    "test",
			Location: "VA1",
			Roles:    []string{"AccountAdmin"},
			Token:    token,
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(auth)
	}
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
	return clc.New(config)
}
