package aa_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/aa"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/stretchr/testify/assert"
)

func TestGetAAPolicy(t *testing.T) {
	assert := assert.New(t)

	name := "12345"
	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.ID)
}

func TestGetAllAAPolicy(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.GetAll()

	assert.Nil(err)
	assert.Equal(2, len(resp.Items))
}

func TestCreateAAPolicy(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	p := aa.Policy{
		Name:     "New AA Policy",
		Location: "va1",
	}

	resp, err := service.Create(p)

	assert.Nil(err)
	assert.Equal(p.Name, resp.Name)
}

func mockStatusAPI() (*httptest.Server, *aa.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/antiAffinityPolicies/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.RequestURI, "/")
		name := parts[len(parts)-1]
		policy := aa.Policy{ID: name}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(policy)

	})
	mux.HandleFunc("/antiAffinityPolicies/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			policies := &aa.Policies{
				Items: []aa.Policy{
					aa.Policy{ID: "123"},
					aa.Policy{ID: "123"},
				},
			}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(policies)
			return
		}

		if r.Method == "POST" {
			policy := &aa.Policy{}
			if err := json.NewDecoder(r.Body).Decode(policy); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			policy.ID = "1235"
			policy.Links = api.Links([]api.Link{api.Link{Rel: "self"}})

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(policy)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
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
	return mockAPI, aa.New(client)
}
