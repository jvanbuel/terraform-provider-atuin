package atuin

import (
	b64 "encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
)

var TEST_API_ENDPOINT = os.Getenv("ATUIN_HOST")

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
	username := "aW0nd3rfulUs3rname"
	password := "password"
	client := NewAtuinClient(TEST_API_ENDPOINT)
	_, err := client.CreateUser(username, password, username+"@example.com")
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}

	err = client.DeleteUser(username, password)
	if err != nil {
		t.Errorf("Error deleting user: %s", err)
	}
}

func TestUpdatePassword(t *testing.T) {
	username := "rincewind"
	password := "swordfish"
	newPassword := "newpassword"
	client := NewAtuinClient(TEST_API_ENDPOINT)

	client.CreateUser(username, password, username+"@example.com")
	err := client.UpdatePassword(username, password, newPassword)
	if err != nil {
		t.Errorf("Error updating password: %s", err)
	}

	client.DeleteUser(username, newPassword)
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

func TestIsValidBip39(t *testing.T) {
	valid := "indoor dish desk flag debris potato excuse depart ticket judge file exit"
	invalid := "er staat een paard in de gang"

	assert.True(t, IsValidBip39(valid))
	assert.False(t, IsValidBip39(invalid))
}
