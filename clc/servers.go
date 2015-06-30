package clc

import "fmt"

type ServerService struct {
	*Client
}

type Server struct {
	Name string
}

func (s *ServerService) Get(name string) (*Server, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.baseURL, s.config.Alias, name)
	server := &Server{}
	err := s.get(url, server)
	return server, err
}
