package clc

import (
	"fmt"
	"time"
)

type Server struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	GroupID        string `json:"groupId"`
	SourceServerID string `json:"sourceServerId"`
	IsManagedOS    bool   `json:"isManagedOS,omitempty"`
	PrimaryDNS     string `json:"primaryDns,omitempty"`
	SecondaryDNS   string `json:"secondaryDns,omitempty"`
	NetworkID      string `json:"networkId,omitempty"`
	IPaddress      string `json:"ipAddress,omitempty"`
	Password       string `json:"password,omitempty"`
	CPU            int    `json:"cpu"`
	MemoryGB       int    `json:"memoryGB"`
	Type           string `json:"type"`
	Storagetype    string `json:"storageType,omitempty"`
	Customfields   []struct {
		ID    string `json:"id"`
		Value string `json:"value"`
	} `json:"customFields,omitempty"`
	Additionaldisks []struct {
		Path   string `json:"path"`
		SizeGB int    `json:"sizeGB"`
		Type   string `json:"type"`
	} `json:"additionalDisks,omitempty"`
	TTL time.Time `json:"ttl,omitempty"`
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
		Customfields []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			Value        string `json:"value"`
			Displayvalue string `json:"displayValue"`
		} `json:"customFields"`
	} `json:"details"`
	Type        string `json:"type"`
	Storagetype string `json:"storageType"`
	ChangeInfo  struct {
		CreatedDate  string `json:"createdDate"`
		CreatedBy    string `json:"createdBy"`
		ModifiedDate string `json:"modifiedDate"`
		ModifiedBy   string `json:"modifiedBy"`
	} `json:"changeInfo"`
	Links []Link `json:"links"`
}

type ServerCreateResponse struct {
	Server   string `json:"server"`
	IsQueued bool   `json:"isQueued"`
	Links    []Link `json:"links"`
}

func (s *ServerCreateResponse) GetStatusId() (string, error) {
	for _, v := range s.Links {
		if v.Rel == "status" {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("No status ID found for server %s", s.Server)
}

type DatacenterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Links []Link `json:"links"`
}

type Link struct {
	Rel   string   `json:"rel,omitempty"`
	Href  string   `json:"href,omitempty"`
	ID    string   `json:"id,omitempty"`
	Verbs []string `json:"verbs,omitempty"`
}
