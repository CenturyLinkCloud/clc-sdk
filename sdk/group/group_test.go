package group_test

import (
	"encoding/json"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/group"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGroup(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Post", "http://localhost/v2/groups/test", mock.Anything, mock.Anything).Return(nil)
	service := group.New(client)

	group := group.Group{
		Name:          "new",
		Description:   "my awesome group",
		Parentgroupid: "12345",
	}
	resp, err := service.Create(group)

	assert.Nil(err)
	assert.Equal(group.Name, resp.Name)
	assert.Equal(1, len(resp.Groups))
	assert.Equal(group.Parentgroupid, resp.Groups[0].ID)
	client.AssertExpectations(t)
}

func TestGetGroup(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/groups/test/67890", mock.Anything).Return(nil)
	service := group.New(client)

	id := "67890"
	resp, err := service.Get(id)

	assert.Nil(err)
	assert.Equal(id, resp.ID)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(url string, resp interface{}) error {
	json.Unmarshal([]byte(mockGroup), resp)
	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) Post(url string, body, resp interface{}) error {
	json.Unmarshal([]byte(mockGroup), resp)
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

const mockGroup = `{"id":"67890","name":"new","description":"my awesome group","locationId":"WA1","type":"default","status":"active","groups":[{"id":"12345","name":"Parent Group Name","description":"The parent group.","locationId":"WA1","type":"default","status":"active","serversCount":0,"groups":[],"links":[{"rel":"createGroup","href":"/v2/groups/acct","verbs":["POST"]},{"rel":"createServer","href":"/v2/servers/acct","verbs":["POST"]},{"rel":"self","href":"/v2/groups/acct/12345","verbs":["GET","PATCH","DELETE"]},{"rel":"parentGroup","href":"/v2/groups/acct/12345","id":"12345"},{"rel":"defaults","href":"/v2/groups/acct/12345/defaults","verbs":["GET","POST"]},{"rel":"billing","href":"/v2/groups/acct/12345/billing"},{"rel":"archiveGroupAction","href":"/v2/groups/acct/12345/archive"},{"rel":"statistics","href":"/v2/groups/acct/12345/statistics"},{"rel":"upcomingScheduledActivities","href":"/v2/groups/acct/12345/upcomingScheduledActivities"},{"rel":"horizontalAutoscalePolicyMapping","href":"/v2/groups/acct/12345/horizontalAutoscalePolicy","verbs":["GET","PUT","DELETE"]},{"rel":"scheduledActivities","href":"/v2/groups/acct/12345/scheduledActivities","verbs":["GET","POST"]}]}],"links":[{"rel":"self","href":"/v2/groups/acct/67890"},{"rel":"parentGroup","href":"/v2/groups/acct/12345","id":"12345"},{"rel":"billing","href":"/v2/groups/acct/67890/billing"},{"rel":"archiveGroupAction","href":"/v2/groups/acct/67890/archive"},{"rel":"statistics","href":"/v2/groups/acct/67890/statistics"},{"rel":"scheduledActivities","href":"/v2/groups/acct/67890/scheduledActivities"}],"changeInfo":{"createdDate":"2012-12-17T01:17:17Z","createdBy":"user@domain.com","modifiedDate":"2014-05-16T23:49:25Z","modifiedBy":"user@domain.com"},"customFields":[]}`
