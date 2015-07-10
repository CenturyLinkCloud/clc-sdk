package clc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

	service := clc.ServerService{Client: client(ms.URL)}
	server, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, server.Name)
}

func TestGetServerByUUID(t *testing.T) {
	assert := assert.New(t)

	name := "5404cf5ece2042dc9f2ac16ab67416bb"
	resource := getServerResource(assert, name)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.ServerService{Client: client(ms.URL)}
	server, err := service.Get(name)

	assert.Nil(err)
	assert.Equal("va1testserver01", server.Name)
}

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	r := createServerResource(assert)
	ms := mockServer(r)
	defer ms.Close()

	service := clc.ServerService{Client: client(ms.URL)}
	server := clc.Server{
		Name:           "va1testserver01",
		CPU:            1,
		MemoryGB:       1,
		GroupID:        "group",
		SourceServerID: "UBUNTU",
		Type:           "standard",
	}
	s, err := service.Create(server, nil)

	assert.Nil(err)
	assert.True(s.IsQueued)
	assert.Equal(server.Name, s.Server)
}

func TestCreateServer_Polling(t *testing.T) {
	assert := assert.New(t)

	r := createServerResource(assert)
	ms := mockServer(r)
	defer ms.Close()

	client := client(ms.URL)
	server := clc.Server{
		Name:           "va1testserver01",
		CPU:            1,
		MemoryGB:       1,
		GroupID:        "group",
		SourceServerID: "UBUNTU",
		Type:           "standard",
	}
	poll := make(chan *clc.StatusResponse, 1)
	_, err := client.Server.Create(server, poll)

	status := <-poll

	assert.Nil(err)
	assert.True(status.Complete())
}

func TestUpdateServer_UpdateCPUAndMemory(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	r := patchServerRequest(assert, name, "cpu", "memory")
	ms := mockServer(r)
	defer ms.Close()

	client := client(ms.URL)
	cpu := clc.ServerCPU(1)
	mem := clc.ServerMemory(1)
	resp, err := client.Server.Update(name, cpu, mem)

	assert.Nil(err)
	assert.Equal(name, resp.Server)
}

func TestDeleteServer(t *testing.T) {
	assert := assert.New(t)

	name := "va1testserver01"
	resource := deleteServerResource(assert, name)
	ms := mockServer(resource)
	defer ms.Close()

	service := clc.ServerService{Client: client(ms.URL)}
	server, err := service.Delete(name)

	assert.Nil(err)
	assert.Equal(name, server.Server)
}

func getServerResource(assert *assert.Assertions, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			assert.Fail("GET server method should be GET", r.Method)
		}

		if r.URL.Path == "/servers/test/"+name && len(r.URL.Query()) == 0 {
			server := &clc.ServerResponse{Name: name}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		if r.URL.Path == "/servers/test/"+name && r.URL.Query().Get("uuid") == "true" {
			server := &clc.ServerResponse{Name: "va1testserver01"}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		assert.Fail("GET server hitting wrong endpoint", r.URL.Path)
	}
}

func patchServerRequest(assert *assert.Assertions, name string, members ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			assert.Fail("PATCH server method should be PATCH", r.Method)
		}

		if r.URL.Path == "/servers/test/"+name {
			updates := make([]clc.ServerUpdate, 0)
			err := json.NewDecoder(r.Body).Decode(&updates)
			if err != nil {
				assert.Fail("body: %s", err)
			}

			for i, v := range updates {
				assert.Equal(members[i], v.Member)
			}

			server := &clc.ServerQueuedResponse{Server: name}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		assert.Fail("PATCH server hitting wrong endpoint", r.URL.Path)
	}
}

func deleteServerResource(assert *assert.Assertions, name string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			assert.Fail("DELETE server method should be DELETE", r.Method)
		}

		if r.URL.Path != "/servers/test/"+name {
			assert.Fail("DELETE server hitting wrong endpoint", r.URL.Path)
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"server":"va1testserver01","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"}]}`)
	}
}

func createServerResource(assert *assert.Assertions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.URL.Path != "/servers/test" {
				assert.Fail("POST server hitting wrong endpoint", r.URL.Path)
			}

			server := &clc.Server{}
			err := json.NewDecoder(r.Body).Decode(server)
			if err != nil {
				assert.Fail("Failed to serialize server", err)
			}

			if !server.Valid() {
				assert.Fail("Server missing required fields", server)
			}

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, `{"server":"va1testserver01","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"},{"rel":"self","href":"/v2/servers/test/12345?uuid=True","id":"12345","verbs":["GET"]}]}`)
		}

		if r.Method == "GET" {
			if !strings.HasPrefix(r.URL.Path, "/operations/test/status/") {
				assert.Fail("Polling hitting wrong endpoint", r.URL.Path)
			}
			status := &clc.StatusResponse{Status: clc.CompleteStatus}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(status)
		}
	}
}
