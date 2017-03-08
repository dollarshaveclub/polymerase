package vaultclient

import (
	"log"
	"os"
	"testing"
)

const (
	testSecretPath = "secret/testing/test_value"
)

var tconfig = VaultConfig{
	Server: os.Getenv("VAULT_ADDR"),
}

func TestVaultAppIDAuth(t *testing.T) {
	vc, err := NewClient(&tconfig)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	err = vc.AppIDAuth("testing-development", "178ae890-4ee1-422d-9877-ed1e784c6adf")
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
}

func TestVaultTokenAuth(t *testing.T) {
	vc, err := NewClient(&tconfig)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	err = vc.TokenAuth(os.Getenv("VAULT_TOKEN"))
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
}

func TestVaultGetValue(t *testing.T) {
	vc, err := NewClient(&tconfig)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	err = vc.AppIDAuth(os.Getenv("VAULT_APP_ID"), os.Getenv("VAULT_USER_ID_PATH"))
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}
	d, err := vc.GetValue(testSecretPath)
	if err != nil {
		log.Fatalf("Error getting value: %v", err)
	}
	log.Printf("Got value: %v", d.(string))
}
