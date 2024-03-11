package atuin

import (
	b64 "encoding/base64"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key, err := GenerateEncryptionKey()
	if err != nil {
		t.Errorf("Error generating key: %s", err)
	}

	_, err = b64.StdEncoding.DecodeString(key)
	if err != nil {
		t.Errorf("Cannot decode key as b64: %s", key)
	}
}

func TestCreateUser(t *testing.T) {
	client := NewAtuinClient(API_ENDPOINT)
	_, err := client.CreateUser("testABDCDEFghijkl", "password", "testABDCDEFghijkl@yahoo.com")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}
}
