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
	return c.DoWithAuth("GET", url, nil, resp)
}

func (c *Client) Post(url string, body, resp interface{}) error {
	return c.DoWithAuth("POST", url, body, resp)
}

func (c *Client) Put(url string, body, resp interface{}) error {
	return c.DoWithAuth("PUT", url, body, resp)
}

func (c *Client) Patch(url string, body, resp interface{}) error {
	return c.DoWithAuth("PATCH", url, body, resp)
}

func (c *Client) Delete(url string, resp interface{}) error {
	return c.DoWithAuth("DELETE", url, nil, resp)
}

func (c *Client) Auth() error {
	url := fmt.Sprintf("%s/authentication/login", c.config.BaseURL)
	body, err := c.serialize(c.config.User)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	return c.Do(req, &c.Token)
}

func (c *Client) Do(req *http.Request, ret interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("http err: [%s]", resp.Status)
	}

	if ret == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(ret)
}

func (c *Client) DoWithAuth(method, url string, body, ret interface{}) error {
	if !c.Token.Valid() {
		err := c.Auth()
		if err != nil {
			return err
		}
	}

	b, err := c.serialize(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.Token.Token)

	return c.Do(req, ret)
}

func (c *Client) serialize(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	return b, err
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
