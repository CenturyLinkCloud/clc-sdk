package clc

import (
	"github.com/mikebeyer/clc-sdk/aa"
	"github.com/mikebeyer/clc-sdk/alert"
	"github.com/mikebeyer/clc-sdk/api"
	"github.com/mikebeyer/clc-sdk/group"
	"github.com/mikebeyer/clc-sdk/lb"
	"github.com/mikebeyer/clc-sdk/server"
	"github.com/mikebeyer/clc-sdk/status"
)

type Client struct {
	client *api.Client

	Server *server.Service
	Status *status.Service
	AA     *aa.Service
	Alert  *alert.Service
	LB     *lb.Service
	Group  *group.Service
}

func New(config api.Config) *Client {
	c := &Client{
		client: api.New(config),
	}

	c.Server = server.New(c.client)
	c.Status = status.New(c.client)
	c.AA = aa.New(c.client)
	c.Alert = alert.New(c.client)
	c.LB = lb.New(c.client)
	c.Group = group.New(c.client)

	return c
}
