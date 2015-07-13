package clc

import "github.com/mikebeyer/clc-sdk/sdk/api"

type Client struct {
	client *api.Client
}

func New(config api.Config) *Client {
	c := &Client{
		client: api.New(config),
	}
	return c
}
