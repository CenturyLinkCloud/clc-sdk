package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mikebeyer/env"
)

func New(config Config) *Client {
	return &Client{
		config: config,
		client: http.DefaultClient,
	}
}

type HTTP interface {
	Get(url string, resp interface{}) error
	Post(url string, body, resp interface{}) error
	Put(url string, body, resp interface{}) error
	Patch(url string, body, resp interface{}) error
	Delete(url string, resp interface{}) error
	Config() *Config
}

type Client struct {
	config Config
	Token  Token

	client *http.Client
}

func (c *Client) Config() *Config {
	return &c.config
}

func (c *Client) Get(url string, resp interface{}) error {
	return c.do("GET", url, nil, resp)
}

func (c *Client) Post(url string, body, resp interface{}) error {
	b, err := c.serialize(body)
	if err != nil {
		return err
	}
	return c.do("POST", url, b, resp)
}

func (c *Client) Put(url string, body, resp interface{}) error {
	b, err := c.serialize(body)
	if err != nil {
		return err
	}
	return c.do("PUT", url, b, resp)
}

func (c *Client) Patch(url string, body, resp interface{}) error {
	b, err := c.serialize(body)
	if err != nil {
		return err
	}
	return c.do("PATCH", url, b, resp)
}

func (c *Client) Delete(url string, resp interface{}) error {
	return c.do("DELETE", url, nil, resp)
}

func (c *Client) serialize(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	return b, err
}

func (c *Client) do(method, url string, body io.Reader, resp interface{}) error {
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

	if res.StatusCode >= 400 {
		return fmt.Errorf("http err: [%s]", res.Status)
	}

	if resp == nil {
		return err
	}
	return json.NewDecoder(res.Body).Decode(resp)
}

func (c *Client) Auth() (string, error) {
	url := fmt.Sprintf("%s/authentication/login", c.config.BaseURL)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(c.config.User)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("http err: [%s]", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.Token); err != nil {
		return "", err
	}

	return c.Token.Token, nil
}

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Username string   `json:"userName"`
	Alias    string   `json:"accountAlias"`
	Location string   `json:"locationAlias"`
	Roles    []string `json:"roles"`
	Token    string   `json:"bearerToken"`
}

// TODO: Add some real validation logic
func (t Token) Valid() bool {
	return t.Token != ""
}

type Update struct {
	Op     string      `json:"op"`
	Member string      `json:"member"`
	Value  interface{} `json:"value"`
}
