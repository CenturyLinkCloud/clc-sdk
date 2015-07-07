package clc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus_Succeeded(t *testing.T) {
	assert := assert.New(t)

	id := "va1-12345"
	status := &clc.StatusResponse{Status: "succeeded"}
	resource := getStatusResource(assert, id, status)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.StatusService{client(ms.URL)}
	status, err := service.Get(id, nil)

	assert.Nil(err)
	assert.True(status.Complete())
}

func TestGetStatus_Poll(t *testing.T) {
	assert := assert.New(t)

	id := "va1-12345"
	ms := mockPollStatus()
	defer ms.Close()

	service := clc.StatusService{client(ms.URL)}
	complete := make(chan *clc.StatusResponse, 1)
	_, err := service.Get(id, complete)

	status := <-complete

	assert.Nil(err)
	assert.Nil(status.Error)
	assert.True(status.Complete())
}

func getStatusResource(assert *assert.Assertions, id string, status *clc.StatusResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path != "/operations/test/status/"+id {
			assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func mockPollStatus() *httptest.Server {
	count := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++

		var status *clc.StatusResponse
		if count <= 1 {
			status = &clc.StatusResponse{Status: "running"}
		} else {
			status = &clc.StatusResponse{Status: clc.CompleteStatus}
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}))
}
