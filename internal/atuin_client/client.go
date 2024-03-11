package atuin

import (
	"bytes"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"net/http"
)

const API_ENDPOINT = "https://api.atuin.sh"

type AtuinClient struct {
	client *http.Client
	host   string
}

func NewAtuinClient(host string) *AtuinClient {
	return &AtuinClient{
		client: &http.Client{},
		host:   host,
	}
}

type Session struct {
	Session string `json:"session"`
}

func (c *AtuinClient) CreateUser(username, password, email string) (string, error) {
	values := map[string]string{"username": username, "password": password, "email": email}

	jsonValue, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", c.host+"/register", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return "", err
	}

	var s Session
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", err
	}

	return s.Session, nil
}

func (c *AtuinClient) UpdatePassword(username, password string) error {
	return nil
}

func (c *AtuinClient) DeleteUser(username, password string) error {
	sessionToken, err := c.Login(username, password)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("DELETE", c.host+"/account", nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Token "+sessionToken)

	_, err = c.client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

func (c *AtuinClient) Login(username, password string) (string, error) {
	values := map[string]string{"username": username, "password": password}

	jsonValue, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", c.host+"/login", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var s Session
	err = json.Unmarshal(body, &s)
	if err != nil {
		return "", err
	}
	return s.Session, nil
}

func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	b64Key := b64.StdEncoding.EncodeToString(key)
	return b64Key, nil
}
