package clc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
