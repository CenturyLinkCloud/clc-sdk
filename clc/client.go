package clc

type Config struct {
	Username string
	Password string
	Alias    string
}

type Client struct {
	config Config
}

func New(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Auth() (string, error) {
	return "", nil
}
