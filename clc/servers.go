package clc

import "fmt"

type ServerService struct {
	*Client
}

func (s *ServerService) Get(name string) (*Server, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.baseURL, s.config.Alias, name)
	server := &Server{}
	err := s.get(url, server)
	return server, err
}

func (s *ServerService) Create(server Server) (*ServerCreate, error) {
	url := fmt.Sprintf("%s/servers/%s", s.baseURL, s.config.Alias)
	resp := &ServerCreate{}
	err := s.post(url, server, resp)
	return resp, err
}
