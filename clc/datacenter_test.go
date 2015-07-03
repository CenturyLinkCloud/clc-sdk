package clc_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetDatacenter(t *testing.T) {
	assert := assert.New(t)

	name := "va1"
	resource := getDatacenterResource(assert, name)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.DatacenterService{client(ms.URL)}
	dc, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, dc.ID)
}

func getDatacenterResource(assert *assert.Assertions, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path != "/datacenters/test/"+name {
			assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
		}

		if r.URL.RawQuery != "groupLinks=true" {
			assert.Fail("GET server providing wrong query params", r.URL.RawQuery)
		}

		dc := &clc.DatacenterResponse{ID: name}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dc)
	}
}
