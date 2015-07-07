package clc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mikebeyer/env"
)

type Config struct {
	User    User
	Alias   string
	BaseURL string
}

func EnvConfig() Config {
	return Config{
		User: User{
			Username: env.MustString("CLC_USERNAME"),
			Password: env.MustString("CLC_PASSWORD"),
		},
		Alias:   env.MustString("CLC_ALIAS"),
		BaseURL: env.String("CLC_BASE_URL", "https://api.ctl.io/v2"),
	}
}

type Client struct {
	config  Config
	client  *http.Client
	baseURL string

	Token Auth

	Server *ServerService
	Status *StatusService
}

func New(config Config) *Client {
	url := config.BaseURL
	if url == "" {
		url = "https://api.ctl.io/v2"
	}
	client := &Client{
		config:  config,
		client:  http.DefaultClient,
		baseURL: url,
	}
	client.Server = &ServerService{client}
	client.Status = &StatusService{client}
	return client
}

func (c *Client) get(url string, resp interface{}) error {
	return c.do("GET", url, nil, resp)
}

func (c *Client) post(url string, body, resp interface{}) error {
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		panic(err)
	}
	return c.do("POST", url, ioutil.NopCloser(b), resp)
}

func (c *Client) do(method, url string, body io.ReadCloser, resp interface{}) error {
	if !c.Token.Valid() {
		token, err := c.Auth()
		if err != nil {
			return err
		}
		c.Token.Token = token
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.Token.Token)
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("http error: %s", res.Status))
	}

	return json.NewDecoder(res.Body).Decode(resp)
}
