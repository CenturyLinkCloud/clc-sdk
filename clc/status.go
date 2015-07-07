package clc

import (
	"fmt"
	"time"
)

type StatusService struct {
	*Client
}

func (s *StatusService) Get(id string, poll chan *StatusResponse) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/operations/%s/status/%s", s.baseURL, s.config.Alias, id)
	status := &StatusResponse{}
	err := s.get(url, status)
	if poll != nil && status.Running() {
		go s.Poll(id, poll)
	}

	return status, err
}

func (s *StatusService) Poll(id string, poll chan *StatusResponse) {
	for {
		state, err := s.Get(id, nil)

		if err != nil {
			poll <- &StatusResponse{Error: err}
			return
		}
		if !state.Running() {
			poll <- state
			return
		}
		time.Sleep(30 * time.Second)
	}
}

type StatusResponse struct {
	Status string `json:"status"`
	Error  error  `json:"-"`
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
