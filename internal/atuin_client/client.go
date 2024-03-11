package atuin

import (
	"crypto/rand"
	b64 "encoding/base64"
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

func (c *AtuinClient) CreatUser(username, password, email string) error {
	return nil
}

func (c *AtuinClient) UpdatePassword(username, password string) error {
	return nil
}

func (c *AtuinClient) DeleteUser(username string) error {
	return nil
}

func (c *AtuinClient) GetSessionToken() (string, error) {
	return "", nil
}

func (c *AtuinClient) GetEncryptionKey() (string, error) {
	return "", nil
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
