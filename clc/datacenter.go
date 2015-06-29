package clc

import "fmt"

func (c *Client) GetDatacenter(name string) (*Datacenter, error) {
	url := fmt.Sprintf("%s/datacenters/%s/%s?groupLinks=true", c.baseURL, c.config.Alias, name)
	dc := &Datacenter{}
	err := c.get(url, dc)
	return dc, err
}

type Datacenter struct {
	ID    string `json:"id"`
	Links []struct {
		Href  string   `json:"href"`
		Rel   string   `json:"rel"`
		Id    string   `json:"id,omitempty"`
		Name  string   `json:"name,omitempty"`
		Verbs []string `json:"verbs,omitempty"`
	} `json:"links"`
	Name string `json:"name"`
}
