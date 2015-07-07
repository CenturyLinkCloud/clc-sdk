package clc

import (
	"errors"
	"fmt"
)

type ServerService struct {
	*Client
}

func (s *ServerService) Get(name string) (*ServerResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.baseURL, s.config.Alias, name)
	server := &ServerResponse{}
	err := s.get(url, server)
	return server, err
}

func (s *ServerService) Create(server Server, poll chan *StatusResponse) (*ServerCreateResponse, error) {
	if !server.Valid() {
		return nil, errors.New("server: server missing required field(s). (Name, CPU, MemoryGB, GroupID, SourceServerID, Type)")
	}

	resp := &ServerCreateResponse{}
	err := s.post(fmt.Sprintf("%s/servers/%s", s.baseURL, s.config.Alias), server, resp)
	if err == nil && resp != nil {
		id, err := resp.GetStatusId()
		if err != nil {
			return resp, err
		}
		go s.Status.Poll(id, poll)
	}

	return resp, err
}
