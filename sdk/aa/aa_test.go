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
	"github.com/stretchr/testify/mock"
)

func TestGetAAPolicy(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	id := "12345"

	resp, err := service.Get(id)

	assert.Nil(err)
	assert.Equal(id, resp.ID)
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

	name := "New AA Policy"
	location := "va1"

	resp, err := service.Create(name, location)

	assert.Nil(err)
	assert.Equal(name, resp.Name)
}

func TestUpdateAAPolicy(t *testing.T) {
	assert := assert.New(t)

	id := "12345"
	name := "My New AA Policy"
	ms, service := mockStatusAPI()
	defer ms.Close()

	resp, err := service.Update(id, name)

	assert.Nil(err)
	assert.Equal(name, resp.Name)
}

func TestDeleteAAPolicy(t *testing.T) {
	assert := assert.New(t)

	ms, service := mockStatusAPI()
	defer ms.Close()

	err := service.Delete("12345")

	assert.Nil(err)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(url string, resp interface{}) error {
	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) Post(url string, body, resp interface{}) error {
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Put(url string, body, resp interface{}) error {
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Patch(url string, body, resp interface{}) error {
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Delete(url string, resp interface{}) error {
	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) GetConfig() *api.Config {
	return &api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: "http://localhost/v2",
	}
}

func mockStatusAPI() (*httptest.Server, *aa.Service) {
	mux := http.NewServeMux()
	mux.HandleFunc("/antiAffinityPolicies/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			parts := strings.Split(r.RequestURI, "/")
			name := parts[len(parts)-1]
			policy := aa.Policy{ID: name}

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(policy)
			return
		}

		if r.Method == "DELETE" {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == "PUT" {
			parts := strings.Split(r.RequestURI, "/")
			id := parts[len(parts)-1]

			policy := &aa.Policy{}
			if err := json.NewDecoder(r.Body).Decode(policy); err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			policy.ID = id
			policy.Location = "va1"
			policy.Links = api.Links([]api.Link{api.Link{Rel: "self"}})

			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(policy)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
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
