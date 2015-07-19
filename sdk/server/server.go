package server

import (
	"fmt"
	"regexp"

	"github.com/mikebeyer/clc-sdk/sdk/api"
)

func New(client api.HTTP) *Service {
	return &Service{
		client: client,
		config: client.Config(),
	}
}

type Service struct {
	client api.HTTP
	config *api.Config
}

func (s *Service) Get(name string) (*Response, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.config.BaseURL, s.config.Alias, name)
	if regexp.MustCompile("^[0-9a-f]{32}$").MatchString(name) {
		url = fmt.Sprintf("%s?uuid=true", url)
	}
	resp := &Response{}
	err := s.client.Get(url, resp)
	return resp, err
}

func (s *Service) Create(server Server) (*QueuedResponse, error) {
	if !server.Valid() {
		return nil, fmt.Errorf("server: server missing required field(s). (Name, CPU, MemoryGB, GroupID, SourceServerID, Type)")
	}

	resp := &QueuedResponse{}
	url := fmt.Sprintf("%s/servers/%s", s.config.BaseURL, s.config.Alias)
	err := s.client.Post(url, server, resp)
	return resp, err
}

func (s *Service) Update(name string, patches ...ServerPatch) (*QueuedResponse, error) {
	resp := &QueuedResponse{}
	url := fmt.Sprintf("%s/servers/%s/%s", s.config.BaseURL, s.config.Alias, name)
	var updates []ServerUpdate
	for _, v := range patches {
		updates = append(updates, v.Serialize())
	}
	err := s.client.Patch(url, updates, resp)
	return resp, err
}

func (s *Service) Delete(name string) (*QueuedResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.config.BaseURL, s.config.Alias, name)
	resp := &QueuedResponse{}
	err := s.client.Delete(url, resp)
	return resp, err
}

func (s *Service) GetPublicIP(name string, ip string) (*PublicIP, error) {
	url := fmt.Sprintf("%s/servers/%s/%s/publicIPAddresses/%s", s.config.BaseURL, s.config.Alias, name, ip)
	resp := &PublicIP{}
	err := s.client.Get(url, resp)
	return resp, err
}

func (s *Service) AddPublicIP(name string, ip PublicIP) (*IPResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s/publicIPAddresses", s.config.BaseURL, s.config.Alias, name)
	resp := &IPResponse{}
	err := s.client.Post(url, ip, resp)
	return resp, err
}

func (s *Service) DeletePublicIP(name, ip string) (*IPResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s/publicIPAddresses/%s", s.config.BaseURL, s.config.Alias, name, ip)
	resp := &IPResponse{}
	err := s.client.Delete(url, resp)
	return resp, err
}

type PublicIP struct {
	InternalIP         string              `json:"internalIPAddress,omitempty"`
	Ports              []Port              `json:"ports,omitempty"`
	SourceRestrictions []SourceRestriction `json:"sourceRestrictions,omitempty"`
}

type Port struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
	PortTo   int    `json:"portTo,omitempty"`
}

type SourceRestriction struct {
	CIDR string `json:"cidr"`
}

type IPResponse struct {
	ID   string `json:"id"`
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type CPU int

func (c CPU) Serialize() ServerUpdate {
	return ServerUpdate{
		Op:     "set",
		Member: "cpu",
		Value:  c,
	}
}

type ServerPatch interface {
	Serialize() ServerUpdate
}

type ServerUpdate struct {
	Op     string `json:"op"`
	Member string `json:"member"`
	Value  interface{}
}

type Server struct {
	Name            string             `json:"name"`
	Description     string             `json:"description,omitempty"`
	GroupID         string             `json:"groupId"`
	SourceServerID  string             `json:"sourceServerId"`
	IsManagedOS     bool               `json:"isManagedOS,omitempty"`
	PrimaryDNS      string             `json:"primaryDns,omitempty"`
	SecondaryDNS    string             `json:"secondaryDns,omitempty"`
	NetworkID       string             `json:"networkId,omitempty"`
	IPaddress       string             `json:"ipAddress,omitempty"`
	Password        string             `json:"password,omitempty"`
	CPU             int                `json:"cpu"`
	MemoryGB        int                `json:"memoryGB"`
	Type            string             `json:"type"`
	Storagetype     string             `json:"storageType,omitempty"`
	Customfields    []api.Customfields `json:"customFields,omitempty"`
	Additionaldisks []struct {
		Path   string `json:"path"`
		SizeGB int    `json:"sizeGB"`
		Type   string `json:"type"`
	} `json:"additionalDisks,omitempty"`
}

func (s *Server) Valid() bool {
	return s.Name != "" && s.CPU != 0 && s.MemoryGB != 0 && s.GroupID != "" && s.SourceServerID != ""
}

type Response struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupID     string `json:"groupId"`
	IsTemplate  bool   `json:"isTemplate"`
	LocationID  string `json:"locationId"`
	OStype      string `json:"osType"`
	Status      string `json:"status"`
	Details     struct {
		IPaddresses []struct {
			Internal string `json:"internal"`
		} `json:"ipAddresses"`
		AlertPolicies []struct {
			ID    string    `json:"id"`
			Name  string    `json:"name"`
			Links api.Links `json:"links"`
		} `json:"alertPolicies"`
		CPU               int    `json:"cpu"`
		Diskcount         int    `json:"diskCount"`
		Hostname          string `json:"hostName"`
		InMaintenanceMode bool   `json:"inMaintenanceMode"`
		MemoryMB          int    `json:"memoryMB"`
		Powerstate        string `json:"powerState"`
		Storagegb         int    `json:"storageGB"`
		Disks             []struct {
			ID             string        `json:"id"`
			SizeGB         int           `json:"sizeGB"`
			PartitionPaths []interface{} `json:"partitionPaths"`
		} `json:"disks"`
		Partitions []struct {
			SizeGB float64 `json:"sizeGB"`
			Path   string  `json:"path"`
		} `json:"partitions"`
		Snapshots []struct {
			Name  string    `json:"name"`
			Links api.Links `json:"links"`
		} `json:"snapshots"`
		Customfields []api.Customfields `json:"customFields,omitempty"`
	} `json:"details"`
	Type        string `json:"type"`
	Storagetype string `json:"storageType"`
	ChangeInfo  struct {
		CreatedDate  string `json:"createdDate"`
		CreatedBy    string `json:"createdBy"`
		ModifiedDate string `json:"modifiedDate"`
		ModifiedBy   string `json:"modifiedBy"`
	} `json:"changeInfo"`
	Links api.Links `json:"links"`
}

type QueuedResponse struct {
	Server   string    `json:"server"`
	IsQueued bool      `json:"isQueued"`
	Links    api.Links `json:"links"`
}

func (q *QueuedResponse) GetStatusID() (bool, string) {
	return q.Links.GetID("status")
}
