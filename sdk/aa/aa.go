package aa

import (
	"fmt"
	"time"

	"github.com/mikebeyer/clc-sdk/sdk/api"
)

func New(client *api.Client) *Service {
	return &Service{
		client:       client,
		PollInterval: 30 * time.Second,
	}
}

type Service struct {
	client *api.Client

	PollInterval time.Duration
}

func (s *Service) Get(id string) (*Policy, error) {
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	policy := &Policy{}
	err := s.client.Get(url, policy)
	return policy, err
}

func (s *Service) GetAll() (*Policies, error) {
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s", s.client.Config.BaseURL, s.client.Config.Alias)
	policies := &Policies{}
	err := s.client.Get(url, policies)
	return policies, err
}

func (s *Service) Create(name, location string) (*Policy, error) {
	policy := &Policy{Name: name, Location: location}
	resp := &Policy{}
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s", s.client.Config.BaseURL, s.client.Config.Alias)
	err := s.client.Post(url, policy, resp)
	return resp, err
}

func (s *Service) Update(id string, name string) (*Policy, error) {
	policy := &Policy{Name: name}
	resp := &Policy{}
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	err := s.client.Put(url, policy, resp)
	return resp, err
}

func (s *Service) Delete(id string) error {
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	err := s.client.Delete(url, nil)
	return err
}

type Policy struct {
	ID       string    `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Location string    `json:"location,omitempty"`
	Links    api.Links `json:"links,omitempty"`
}

type Policies struct {
	Items []Policy  `json:"items,omitempty"`
	Links api.Links `json:"links,omitempty"`
}
