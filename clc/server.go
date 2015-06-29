package clc

import "fmt"

func (c *Client) GetServer(name string) (*ServerResponse, error) {
	url := fmt.Sprintf("%s/servers/%s/%s", c.baseURL, c.config.Alias, name)
	server := &ServerResponse{}
	err := c.get(url, server)
	return server, err
}

type Server struct {
	CPU             int    `json:"cpu"`
	Description     string `json:"description"`
	GroupID         string `json:"groupId"`
	IPAddress       string `json:"ipAddress"`
	IsManagedOS     bool   `json:"isManagedOS"`
	MemoryGB        int    `json:"memoryGB"`
	Name            string `json:"name"`
	NetworkID       string `json:"networkId"`
	Password        string `json:"password"`
	PrimaryDNS      string `json:"primaryDns"`
	SecondaryDNS    string `json:"secondaryDns"`
	SourceServerID  string `json:"sourceServerId"`
	StorageType     string `json:"storageType"`
	Type            string `json:"type"`
	CustomFields    map[string]string
	AdditionalDisks []struct {
		Path   string `json:"path"`
		SizeGB int    `json:"sizeGB"`
		Type   string `json:"type"`
	} `json:"additionalDisks"`
}

type ServerResponse struct {
	ChangeInfo struct {
		CreatedBy    string `json:"createdBy"`
		CreatedDate  string `json:"createdDate"`
		ModifiedBy   string `json:"modifiedBy"`
		ModifiedDate string `json:"modifiedDate"`
	} `json:"changeInfo"`
	Description string `json:"description"`
	Details     struct {
		AlertPolicies []interface{} `json:"alertPolicies"`
		CPU           int           `json:"cpu"`
		CustomFields  []interface{} `json:"customFields"`
		DiskCount     int           `json:"diskCount"`
		Disks         []struct {
			ID             string        `json:"id"`
			PartitionPaths []interface{} `json:"partitionPaths"`
			SizeGB         int           `json:"sizeGB"`
		} `json:"disks"`
		HostName          string `json:"hostName"`
		InMaintenanceMode bool   `json:"inMaintenanceMode"`
		IPAddresses       []struct {
			Internal string `json:"internal"`
		} `json:"ipAddresses"`
		MemoryMB   int `json:"memoryMB"`
		Partitions []struct {
			Path   string  `json:"path"`
			SizeGB float64 `json:"sizeGB"`
		} `json:"partitions"`
		PowerState string        `json:"powerState"`
		Snapshots  []interface{} `json:"snapshots"`
		StorageGB  int           `json:"storageGB"`
	} `json:"details"`
	GroupID    string `json:"groupId"`
	ID         string `json:"id"`
	IsTemplate bool   `json:"isTemplate"`
	Links      []struct {
		Href  string   `json:"href"`
		ID    string   `json:"id"`
		Rel   string   `json:"rel"`
		Verbs []string `json:"verbs"`
	} `json:"links"`
	LocationID  string `json:"locationId"`
	Name        string `json:"name"`
	Os          string `json:"os"`
	OsType      string `json:"osType"`
	Status      string `json:"status"`
	StorageType string `json:"storageType"`
	Type        string `json:"type"`
}
