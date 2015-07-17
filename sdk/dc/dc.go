package dc

import (
	"fmt"
	"time"

	"github.com/mikebeyer/clc-sdk/sdk/api"
)

func New(client api.HTTP) *Service {
	return &Service{
		client:       client,
		PollInterval: 30 * time.Second,
	}
}

type Service struct {
	client api.HTTP

	PollInterval time.Duration
}

func (s *Service) Get(id string) (*Response, error) {
	config := s.client.GetConfig()
	url := fmt.Sprintf("%s/datacenters/%s/%s?groupLinks=true", config.BaseURL, config.Alias, id)
	dc := &Response{}
	err := s.client.Get(url, dc)
	return dc, err
}

func (s *Service) GetAll() ([]*Response, error) {
	config := s.client.GetConfig()
	url := fmt.Sprintf("%s/datacenters/%s", config.BaseURL, config.Alias)
	dcs := make([]*Response, 0)
	err := s.client.Get(url, &dcs)
	return dcs, err
}

type Response struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Links api.Links `json:"links"`
}
