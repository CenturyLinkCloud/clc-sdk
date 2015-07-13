package status

import (
	"fmt"

	"github.com/mikebeyer/clc-sdk/sdk/api"
)

func New(client *api.Client) *Service {
	return &Service{client}
}

type Service struct {
	client *api.Client
}

func (s *Service) Get(id string, poll chan *Response) (*Response, error) {
	url := fmt.Sprintf("%s/operations/%s/status/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	status := &Response{}
	err := s.client.Get(url, status)

	return status, err
}

type Response struct {
	Status string `json:"status"`
	Error  error  `json:"-"`
}

func (s *Response) Complete() bool {
	return s.Status == Complete
}

func (s *Response) Failed() bool {
	return s.Status == Failed
}

func (s *Response) Running() bool {
	return !s.Complete() && !s.Failed()
}

const (
	Complete = "succeeded"
	Failed   = "failed"
)
