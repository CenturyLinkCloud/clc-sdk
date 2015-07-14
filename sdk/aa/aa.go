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

func (s *Service) Get(id string) (*Response, error) {
	url := fmt.Sprintf("%s/antiAffinityPolicies/%s/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	status := &Response{}
	err := s.client.Get(url, status)
	return status, err
}

type Response struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Links    api.Links `json:"links"`
}
