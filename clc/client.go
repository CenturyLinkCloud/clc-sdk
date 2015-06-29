package clc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	config  *Config
	client  *http.Client
	token   Token
	baseURL string
}

func New(config *Config) *Client {
	return &Client{
		config:  config,
		client:  http.DefaultClient,
		baseURL: "https://api.ctl.io/v2",
	}
}

func (c *Client) get(url string, resp interface{}) error {
	return c.do("get", url, nil, resp)
}

func (c *Client) do(method, url string, body io.Reader, resp interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(resp)
}

func (c *Client) auth() (string, error) {
	url := fmt.Sprintf("%s/authentication/login", c.baseURL)
	body := []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, c.config.Name, c.config.Password))
	req, err := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if err := json.NewDecoder(resp.Body).Decode(&c.token); err != nil {
		return "", err
	}
	return c.token.Token, nil
}

type Token struct {
	Token string `json:"bearerToken"`
}

func (t Token) Bearer() string {
	return fmt.Sprintf("Bearer %s", t.Token)
}

func (t Token) Exp() bool {
	return t.Token != ""
}

type Config struct {
	Name     string
	Password string
	Alias    string
}

func EnvConfig() (*Config, error) {
	user := os.Getenv("CLC_USERNAME")
	if user == "" {
		return nil, errors.New("Please set CLC_USERNAME")
	}
	pw := os.Getenv("CLC_PASSWORD")
	if pw == "" {
		return nil, errors.New("Please set CLC_PASSWORD")
	}

	return &Config{
		Name:     user,
		Password: pw,
	}, nil
}
