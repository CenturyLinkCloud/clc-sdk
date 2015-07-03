package clc

import "fmt"

type DatacenterService struct {
	*Client
}

func (s *DatacenterService) Get(name string) (*DatacenterResponse, error) {
	url := fmt.Sprintf("%s/datacenters/%s/%s?groupLinks=true", s.baseURL, s.config.Alias, name)
	server := &DatacenterResponse{}
	err := s.get(url, server)
	return server, err
}
