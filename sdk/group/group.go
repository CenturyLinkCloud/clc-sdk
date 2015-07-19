package group

import (
	"fmt"
	"time"

	"github.com/mikebeyer/clc-sdk/sdk/api"
)

func New(client api.HTTP) *Service {
	return &Service{
		client: client,
		config: client.Config(),
	}
}

type Service struct {
	client api.HTTP
	config *api.Config
}

func (s *Service) Get(id string) (*Response, error) {
	url := fmt.Sprintf("%s/groups/%s/%s", s.config.BaseURL, s.config.Alias, id)
	resp := &Response{}
	err := s.client.Get(url, resp)
	return resp, err
}

func (s *Service) Create(group Group) (*Response, error) {
	resp := &Response{}
	url := fmt.Sprintf("%s/groups/%s", s.config.BaseURL, s.config.Alias)
	err := s.client.Post(url, group, resp)
	return resp, err
}

type Response struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Locationid  string    `json:"locationId"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Links       api.Links `json:"links"`
	Groups      []Groups  `json:"groups"`
	Changeinfo  struct {
		Createddate  time.Time `json:"createdDate"`
		Createdby    string    `json:"createdBy"`
		Modifieddate time.Time `json:"modifiedDate"`
		Modifiedby   string    `json:"modifiedBy"`
	} `json:"changeInfo"`
	Customfields []api.Customfields `json:"customFields"`
}
type Group struct {
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Parentgroupid string             `json:"parentGroupId"`
	Customfields  []api.Customfields `json:"customFields"`
}

type Groups struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Locationid   string    `json:"locationId"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	Serverscount int       `json:"serversCount"`
	Groups       []Groups  `json:"groups"`
	Links        api.Links `json:"links"`
}
