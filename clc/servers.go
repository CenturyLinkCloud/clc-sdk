package clc

import "fmt"

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
	url := fmt.Sprintf("%s/servers/%s", s.baseURL, s.config.Alias)
	resp := &ServerCreateResponse{}
	err := s.post(url, server, resp)
	return resp, err
}
