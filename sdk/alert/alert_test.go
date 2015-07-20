package alert_test

import (
	"encoding/json"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/alert"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAlertPolicy(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/alertPolicies/test/12345", mock.Anything, mock.Anything).Return(nil)
	service := alert.New(client)

	id := "12345"
	resp, err := service.Get(id)

	assert.Nil(err)
	assert.Equal(id, resp.ID)
	client.AssertExpectations(t)
}

func TestCreateAlertPolicy(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Post", "http://localhost/v2/alertPolicies/test", mock.Anything, mock.Anything).Return(nil)
	service := alert.New(client)

	a := alert.Alert{
		Name: "new alert",
		Actions: []alert.Action{
			alert.Action{
				Action: "email",
				Setting: alert.Setting{
					[]string{"user@company.com"},
				},
			},
		},
		Triggers: []alert.Trigger{
			alert.Trigger{
				Metric:    "disk",
				Duration:  "00:05:00",
				Threshold: 80.0,
			},
		},
	}
	resp, err := service.Create(a)

	assert.Nil(err)
	assert.Equal(a.Name, resp.Name)
	assert.Equal(a.Actions, resp.Actions)
	assert.Equal(a.Triggers, resp.Triggers)
	client.AssertExpectations(t)
}

func TestUpdateAlertPolicy(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Put", "http://localhost/v2/alertPolicies/test/12345", mock.Anything, mock.Anything).Return(nil)
	service := alert.New(client)

	a := alert.Alert{
		Name: "update alert",
	}
	id := "12345"
	resp, err := service.Update(id, a)

	assert.Nil(err)
	assert.Equal(a.Name, resp.Name)
	assert.Equal("12345", resp.ID)
	client.AssertExpectations(t)
}

func TestDeleteAlertPolicy(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Delete", "http://localhost/v2/alertPolicies/test/12345", mock.Anything).Return(nil)
	service := alert.New(client)

	err := service.Delete("12345")

	assert.Nil(err)
	client.AssertExpectations(t)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(url string, resp interface{}) error {
	json.Unmarshal([]byte(`{"id":"12345","name":"new alert","actions":[{"action":"email","settings":{"recipients":["user@company.com"]}}],"links":[{"rel":"self","href":"/v2/alertPolicies/test/12345","verbs":["GET","DELETE","PUT"]}],"triggers":[{"metric":"disk","duration":"00:05:00","threshold":80.0}]}`), resp)
	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) Post(url string, body, resp interface{}) error {
	json.Unmarshal([]byte(`{"id":"12345","name":"new alert","actions":[{"action":"email","settings":{"recipients":["user@company.com"]}}],"links":[{"rel":"self","href":"/v2/alertPolicies/test/12345","verbs":["GET","DELETE","PUT"]}],"triggers":[{"metric":"disk","duration":"00:05:00","threshold":80.0}]}`), resp)
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Put(url string, body, resp interface{}) error {
	json.Unmarshal([]byte(`{"id":"12345","name":"update alert","actions":[{"action":"email","settings":{"recipients":["user@company.com"]}}],"links":[{"rel":"self","href":"/v2/alertPolicies/test/12345","verbs":["GET","DELETE","PUT"]}],"triggers":[{"metric":"disk","duration":"00:05:00","threshold":80.0}]}`), resp)
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
