package clc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Config struct {
	User    User
	Alias   string
	BaseURL string
}

type Client struct {
	config  Config
	client  *http.Client
	baseURL string
}

func New(config Config) *Client {
	url := config.BaseURL
	if url == "" {
		url = "https://api.ctl.io/v2"
	}
	return &Client{
		config:  config,
		client:  http.DefaultClient,
		baseURL: url,
	}
}

func (c *Client) Auth() (string, error) {
	url := fmt.Sprintf("%s/authentication/login", c.baseURL)
	b, err := json.Marshal(&c.config.User)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", ioutil.NopCloser(bytes.NewReader(b)))
	if err != nil {
		return "", err
	}

	auth := &Auth{}
	if err := json.NewDecoder(resp.Body).Decode(auth); err != nil {
		return "", err
	}

	return auth.Token, nil
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Username string   `json:"userName"`
	Alias    string   `json:"accountAlias"`
	Location string   `json:"locationAlias"`
	Roles    []string `json:"roles"`
	Token    string   `json:"bearerToken"`
}
