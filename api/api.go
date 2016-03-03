package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/mikebeyer/env"
)

var debug = os.Getenv("DEBUG") != ""

const baseUriDefault = "https://api.ctl.io/v2"
const userAgentDefault = "CenturyLinkCloud/clc-sdk"

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

	return c.Do(req, &c.Token)
}

func (c *Client) Do(req *http.Request, ret interface{}) error {
	if debug {
		v, _ := httputil.DumpRequest(req, true)
		log.Println(string(v))
	}

	req.Header.Add("User-Agent", c.config.UserAgent)
	req.Header.Add("Accept", "application/json")
	if req.Body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if debug && resp != nil {
		v, _ := httputil.DumpResponse(resp, true)
		log.Println(string(v))
	}
	if resp.StatusCode >= 400 {
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("http err: [%s] - %s", resp.Status, b)
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
	User      User     `json:"user"`
	Alias     string   `json:"alias"`
	BaseURL   *url.URL `json:"-"`
	UserAgent string   `json:"agent,omitempty"`
}

func (c Config) Valid() bool {
	return c.User.Username != "" && c.User.Password != "" && c.Alias != "" && c.BaseURL != nil
}

func EnvConfig() (Config, error) {
	user := os.Getenv("CLC_USERNAME")
	pass := os.Getenv("CLC_PASSWORD")
	alias := os.Getenv("CLC_ALIAS")
	base := env.String("CLC_BASE_URL", baseUriDefault)

	config, err := NewConfig(user, pass, alias, base)
	if ua := os.Getenv("CLC_USER_AGENT"); ua != "" {
		config.UserAgent = ua
	}

	if err != nil {
		return config, err
	}

	if !config.Valid() {
		return config, fmt.Errorf("missing environment variables [%s]", config)
	}
	return config, nil
}

func NewConfig(username, password, alias string, uri string) (Config, error) {
	if uri == "" {
		uri = baseUriDefault
	}

	u, err := url.Parse(uri)

	return Config{
		User: User{
			Username: username,
			Password: password,
		},
		Alias:     alias,
		BaseURL:   u,
		UserAgent: userAgentDefault,
	}, err
}

func FileConfig(file string) (Config, error) {
	config := Config{}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(b, &config)

	u, err := url.Parse(baseUriDefault)
	config.BaseURL = u
	return config, err
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
