package atuin

import (
	"bytes"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tyler-smith/go-bip39"
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

func (c *AtuinClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/json")
	return c.client.Do(req)
}

func (c *AtuinClient) CreateUser(username, password, email string) (string, error) {
	values := map[string]string{"username": username, "password": password, "email": email}

	jsonValue, _ := json.Marshal(values)

	request, err := http.NewRequest("POST", c.host+"/register", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	resp, err := c.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error creating user: %s", resp.Status)
	}

	defer resp.Body.Close()

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

	resp, err := c.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating user: %s", resp.Status)
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

	resp, err := c.Do(request)
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

func ConvertEncryptionKeyToBip39(key string) (string, error) {
	decodedKey, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(decodedKey)
}

func IsValidBip39(key string) bool {
	_, err := bip39.EntropyFromMnemonic(key)
	return err == nil
}
