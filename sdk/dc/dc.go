package dc

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

func (s *Service) Get(id string) (*Response, error) {
	url := fmt.Sprintf("%s/datacenters/%s/%s?groupLinks=true", s.config.BaseURL, s.config.Alias, id)
	dc := &Response{}
	err := s.client.Get(url, dc)
	return dc, err
}

func (s *Service) GetAll() ([]*Response, error) {
	url := fmt.Sprintf("%s/datacenters/%s", s.config.BaseURL, s.config.Alias)
	dcs := make([]*Response, 0)
	err := s.client.Get(url, &dcs)
	return dcs, err
}

type Response struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Links api.Links `json:"links"`
}
