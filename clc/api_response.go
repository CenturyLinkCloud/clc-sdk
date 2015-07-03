package clc

type Server struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupID     string `json:"groupId"`
	Istemplate  bool   `json:"isTemplate"`
	LocationID  string `json:"locationId"`
	Ostype      string `json:"osType"`
	Status      string `json:"status"`
	Details     struct {
		Ipaddresses []struct {
			Internal string `json:"internal"`
		} `json:"ipAddresses"`
		Alertpolicies []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Links []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
		} `json:"alertPolicies"`
		CPU               int    `json:"cpu"`
		Diskcount         int    `json:"diskCount"`
		Hostname          string `json:"hostName"`
		Inmaintenancemode bool   `json:"inMaintenanceMode"`
		Memorymb          int    `json:"memoryMB"`
		Powerstate        string `json:"powerState"`
		Storagegb         int    `json:"storageGB"`
		Disks             []struct {
			ID             string        `json:"id"`
			Sizegb         int           `json:"sizeGB"`
			Partitionpaths []interface{} `json:"partitionPaths"`
		} `json:"disks"`
		Partitions []struct {
			Sizegb float64 `json:"sizeGB"`
			Path   string  `json:"path"`
		} `json:"partitions"`
		Snapshots []struct {
			Name  string `json:"name"`
			Links []struct {
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"links"`
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
	Changeinfo  struct {
		Createddate  string `json:"createdDate"`
		Createdby    string `json:"createdBy"`
		Modifieddate string `json:"modifiedDate"`
		Modifiedby   string `json:"modifiedBy"`
	} `json:"changeInfo"`
	Links []Link `json:"links"`
}

type ServerCreate struct {
	Server   string `json:"server"`
	IsQueued bool   `json:"isQueued"`
	Links    []Link `json:"links"`
}

type Link struct {
	Rel   string   `json:"rel"`
	Href  string   `json:"href"`
	ID    string   `json:"id"`
	Verbs []string `json:"verbs,omitempty"`
}
