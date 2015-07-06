package clc_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetStatus_Succeeded(t *testing.T) {
	assert := assert.New(t)

	id := "va1-12345"
	resource := getStatusResource(assert, id, "succeeded")
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.StatusService{client(ms.URL)}
	status, err := service.Get(id)

	assert.Nil(err)
	assert.True(status.Complete())
}

func getStatusResource(assert *assert.Assertions, id string, status string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path != "/operations/test/status/"+id {
			assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
		}

		status := clc.StatusResponse{status}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
