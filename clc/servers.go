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

func (s *ServerService) Create(server Server) (*ServerCreateResponse, error) {
	if !server.Valid() {
		return nil, errors.New("server: server missing required field(s). (Name, CPU, MemoryGB, GroupID, SourceServerID, Type)")
	}
	url := fmt.Sprintf("%s/servers/%s", s.baseURL, s.config.Alias)
	resp := &ServerCreateResponse{}
	err := s.post(url, server, resp)
	return resp, err
}
