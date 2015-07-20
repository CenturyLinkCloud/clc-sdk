package alert

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

func (s *Service) Create(alert Alert) (*Alert, error) {
	url := fmt.Sprintf("%s/alertPolicies/%s", s.config.BaseURL, s.config.Alias)
	resp := &Alert{}
	err := s.client.Post(url, alert, resp)
	return resp, err
}

type Alert struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Actions  []Action  `json:"actions"`
	Triggers []Trigger `json:"triggers"`
	Links    api.Links `json:"links"`
}

type Action struct {
	Action  string  `json:"action"`
	Setting Setting `json:"settings"`
}

type Setting struct {
	Recipients []string `json:"recipients"`
}

type Trigger struct {
	Metric    string  `json:"metric"`
	Duration  string  `json:"duration"`
	Threshold float64 `json:"threshold"`
}
