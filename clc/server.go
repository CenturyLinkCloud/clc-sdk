package clc

import (
	"errors"
	"fmt"
	"regexp"
)

type ServerService struct {
	*Client
}

func (s *ServerService) Get(name string) (*ServerResponse, error) {
	var query string
	var uuidRegex = regexp.MustCompile("[0-9a-f]{8}[0-9a-f]{4}[0-9a-f]{4}[0-9a-f]{4}[0-9a-f]{12}")

	if uuidRegex.MatchString(name) {
		query = "?uuid=true"
	}
	url := fmt.Sprintf("%s/servers/%s/%s%s", s.baseURL, s.config.Alias, name, query)
	server := &ServerResponse{}
	err := s.get(url, server)
	return server, err
}

func (s *ServerService) Create(server Server, poll chan *StatusResponse) (*ServerQueuedResponse, error) {
	if !server.Valid() {
		return nil, errors.New("server: server missing required field(s). (Name, CPU, MemoryGB, GroupID, SourceServerID, Type)")
	}

	if server.Type == "hyperscale" {
		server.Storagetype = "hyperscale"
	}

	resp := &ServerQueuedResponse{}
	err := s.post(fmt.Sprintf("%s/servers/%s", s.baseURL, s.config.Alias), server, resp)
	if err == nil && poll != nil {
		ok, id := resp.Links.GetID("status")
		if !ok {
			return resp, fmt.Errorf("No status ID avaiable to poll for server: %s", resp.Server)
		}
		go s.Status.Poll(id, poll)
	}

	return resp, err
}

func (s *ServerService) Update(name string, patches ...ServerPatch) (*ServerQueuedResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.baseURL, s.config.Alias, name)
	server := &ServerQueuedResponse{}
	var updates []ServerUpdate
	for _, v := range patches {
		m, val := v.Value()
		updates = append(updates, ServerUpdate{Op: "set", Member: m, Value: val})
	}
	err := s.patch(url, updates, server)
	return server, err
}

func (s *ServerService) Delete(name string) (*ServerQueuedResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", s.baseURL, s.config.Alias, name)
	server := &ServerQueuedResponse{}
	err := s.delete(url, server)
	return server, err
}

type ServerPatch interface {
	Value() (string, interface{})
}

type ServerUpdate struct {
	Op     string      `json:"op"`
	Member string      `json:"member"`
	Value  interface{} `json:"value"`
}

type ServerCPU int

func (c ServerCPU) Value() (string, interface{}) {
	return "cpu", c
}

type ServerMemory int

func (m ServerMemory) Value() (string, interface{}) {
	return "memory", m
}

type ServerDescription string

func (d ServerDescription) Value() (string, interface{}) {
	return "description", d
}

type Server struct {
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	GroupID         string         `json:"groupId"`
	SourceServerID  string         `json:"sourceServerId"`
	IsManagedOS     bool           `json:"isManagedOS,omitempty"`
	PrimaryDNS      string         `json:"primaryDns,omitempty"`
	SecondaryDNS    string         `json:"secondaryDns,omitempty"`
	NetworkID       string         `json:"networkId,omitempty"`
	IPaddress       string         `json:"ipAddress,omitempty"`
	Password        string         `json:"password,omitempty"`
	CPU             int            `json:"cpu"`
	MemoryGB        int            `json:"memoryGB"`
	Type            string         `json:"type"`
	Storagetype     string         `json:"storageType,omitempty"`
	Customfields    []Customfields `json:"customFields,omitempty"`
	Additionaldisks []struct {
		Path   string `json:"path"`
		SizeGB int    `json:"sizeGB"`
		Type   string `json:"type"`
	} `json:"additionalDisks,omitempty"`
}

func (s *Server) Valid() bool {
	return s.Name != "" && s.CPU != 0 && s.MemoryGB != 0 && s.GroupID != "" && s.SourceServerID != ""
}

type ServerResponse struct {
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
			ID    string `json:"id"`
			Name  string `json:"name"`
			Links []Link `json:"links"`
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
			Name  string `json:"name"`
			Links []Link `json:"links"`
		} `json:"snapshots"`
		Customfields []Customfields `json:"customFields,omitempty"`
	} `json:"details"`
	Type        string `json:"type"`
	Storagetype string `json:"storageType"`
	ChangeInfo  struct {
		CreatedDate  string `json:"createdDate"`
		CreatedBy    string `json:"createdBy"`
		ModifiedDate string `json:"modifiedDate"`
		ModifiedBy   string `json:"modifiedBy"`
	} `json:"changeInfo"`
	Links Links `json:"links"`
}

type Customfields struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Value        string `json:"value,omitempty"`
	Displayvalue string `json:"displayValue,omitempty"`
}

type ServerQueuedResponse struct {
	Server   string `json:"server"`
	IsQueued bool   `json:"isQueued"`
	Links    Links  `json:"links"`
}
