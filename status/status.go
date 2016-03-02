package status

import (
	"fmt"
	"time"

	"github.com/CenturyLinkCloud/clc-sdk/api"
)

func New(client api.HTTP) *Service {
	return &Service{
		client:       client,
		config:       client.Config(),
		PollInterval: 30 * time.Second,
	}
}

type Service struct {
	client api.HTTP
	config *api.Config

	PollInterval time.Duration
}

func (s *Service) Get(id string) (*Response, error) {
	url := fmt.Sprintf("%s/operations/%s/status/%s", s.config.BaseURL, s.config.Alias, id)
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

type Status struct {
	ID   string `json:"id"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
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
	return !s.Complete() && !s.Failed() && s.Status != ""
}

const (
	Complete = "succeeded"
	Failed   = "failed"
)

type QueuedResponse struct {
	Server   string    `json:"server,omitempty"`
	IsQueued bool      `json:"isQueued,omitempty"`
	Links    api.Links `json:"links,omitempty"`
	Error    string    `json:"errorMessage,omitempty"`
}

func (q *QueuedResponse) GetStatusID() (bool, string) {
	return q.Links.GetID("status")
}
