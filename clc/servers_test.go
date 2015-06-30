package clc_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/stretchr/testify/assert"
)

func TestGetServer(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	resource := getServerResource(assert, name)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.ServerService{client(ms.URL)}
	server, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, server.Name)
}

func getServerResource(assert *assert.Assertions, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path != "/servers/test/"+name {
			assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
		}

		server := &clc.Server{Name: name}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(server)
	}
}
