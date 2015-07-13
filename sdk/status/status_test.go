package status_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/status"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Get("12345", nil)

	assert.Nil(err)
	assert.True(resp.Complete())
}

func mockStatusAPI() (*httptest.Server, *status.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/operations/test/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"succeeded"}`)
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

	return mockAPI, status.New(client)
}
