package status

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
	url := fmt.Sprintf("%s/operations/%s/status/%s", s.client.Config.BaseURL, s.client.Config.Alias, id)
	status := &Response{}
	err := s.client.Get(url, status)
	return status, err
}

func (s *Service) Poll(id string, poll chan *Response) error {
	for {
		status, err := s.Get(id)
		if err != nil {
			return err
		}

		if !status.Running() {
			poll <- status
			return nil
		}
		time.Sleep(s.PollInterval)
	}
}

type Response struct {
	Status string `json:"status"`
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
