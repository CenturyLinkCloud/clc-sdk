package status_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/mikebeyer/clc-sdk/api"
	"github.com/mikebeyer/clc-sdk/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStatus(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/operations/test/status/12345", mock.Anything).Return(nil)
	service := status.New(client)
	resp, err := service.Get("12345")

	assert.Nil(err)
	assert.True(resp.Running())
	client.AssertExpectations(t)
}

func TestPollStatus(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/operations/test/status/12345", mock.Anything).Return(nil)
	service := status.New(client)
	service.PollInterval = 1 * time.Microsecond

	poll := make(chan *status.Response, 1)
	err := service.Poll("12345", poll)

	status := <-poll

	assert.Nil(err)
	assert.True(status.Complete())
	client.AssertExpectations(t)
}

func TestPollStatus_ErrorGettingStatus(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/operations/test/status/12345", mock.Anything).Return(errors.New(""))
	service := status.New(client)
	service.PollInterval = 1 * time.Microsecond

	err := service.Poll("12345", make(chan *status.Response, 1))

	assert.NotNil(err)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	mock.Mock

	count int
}

func (m *MockClient) Get(url string, resp interface{}) error {
	if m.count <= 1 {
		json.Unmarshal([]byte(`{"status":"executing"}`), resp)
	} else {
		json.Unmarshal([]byte(`{"status":"succeeded"}`), resp)
	}
	m.count++
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

func (m *MockClient) Config() *api.Config {
	return &api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: "http://localhost/v2",
	}
}
