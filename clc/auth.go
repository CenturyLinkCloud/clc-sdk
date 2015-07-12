package clc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) Auth() (string, error) {
	url := fmt.Sprintf("%s/authentication/login", c.baseURL)
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

func (a Auth) Valid() bool {
	return a.Token != ""
}
