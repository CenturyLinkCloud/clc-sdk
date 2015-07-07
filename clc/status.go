package clc

import "fmt"

type StatusService struct {
	*Client
}

func (s *StatusService) Get(id string) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/operations/%s/status/%s", s.baseURL, s.config.Alias, id)
	server := &StatusResponse{}
	err := s.get(url, server)
	return server, err
}

type StatusResponse struct {
	Status string `json:"status"`
}

func (s *StatusResponse) Complete() bool {
	return s.Status == CompleteStatus
}

func (s *StatusResponse) Failed() bool {
	return s.Status == FailedStatus
}

func (s *StatusResponse) Running() bool {
	return !s.Complete() && !s.Failed()
}

const (
	CompleteStatus = "succeeded"
	FailedStatus   = "failed"
)
