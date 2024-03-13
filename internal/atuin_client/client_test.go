package atuin

import (
	b64 "encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
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

func TestCreateAndDeleteUser(t *testing.T) {
	username := "randomusernameABCDEF"
	password := "password"
	client := NewAtuinClient(API_ENDPOINT)
	_, err := client.CreateUser(username, password, username+"@example.com")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	err = client.DeleteUser(username, password)
	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}
}

func TestConvertKeyToBip39(t *testing.T) {
	randomKey, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatal(err)
	}

	bip39Key, err := ConvertEncryptionKeyToBip39(randomKey)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := bip39.EntropyFromMnemonic(bip39Key)
	if err != nil {
		t.Fatal(err)
	}

	decodedb64 := b64.StdEncoding.EncodeToString(decoded)

	assert.Equal(t, decodedb64, randomKey)
}
