package lb

import (
	"fmt"

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

func (s *Service) Create(dc string, lb LoadBalancer) (*LoadBalancer, error) {
	url := fmt.Sprintf("%s/sharedLoadBalancers/%s/%s", s.config.BaseURL, s.config.Alias, dc)
	resp := &LoadBalancer{}
	err := s.client.Post(url, lb, resp)
	return resp, err
}

type LoadBalancer struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IPaddress   string    `json:"ipAddress,omitempty"`
	Status      string    `json:"status,omitempty"`
	Pools       []Pool    `json:"pools,omitempty"`
	Links       api.Links `json:"links,omitempty"`
}

type Pool struct {
}
