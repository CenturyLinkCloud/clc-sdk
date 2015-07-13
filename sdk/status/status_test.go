package status_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/status"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Get("12345")

	assert.Nil(err)
	assert.True(resp.Complete())
}

func TestGetStatus_Polling(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	service.PollInterval = 1 * time.Microsecond
	defer ms.Close()

	poll := make(chan *status.Response, 1)
	err := service.Poll("poll", poll)

	status := <-poll

	assert.Nil(err)
	assert.True(status.Complete())
}

func mockStatusAPI() (*httptest.Server, *status.Service) {
	count := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/operations/test/status/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.RequestURI == "/operations/test/status/poll" {
			w.Header().Add("Content-Type", "application/json")
			if count <= 1 {
				fmt.Fprintf(w, `{"status":"running"}`)
			} else {
				fmt.Fprintf(w, `{"status":"succeeded"}`)
			}
			count++
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
