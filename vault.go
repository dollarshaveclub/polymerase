package main

import "github.com/dollarshaveclub/polymerase/pkg/vaultclient"

// Vault is a simple interface for a vault client
type Vault interface {
	GetStringValue(string) (string, error)
}

// AuthenticatedVaultClient creates and authenicates a vault client using the given config
func AuthenticatedVaultClient(config Config) (Vault, error) {

	v, err := vaultclient.NewClient(&vaultclient.VaultConfig{Server: config.VaultAddr})
	if err != nil {
		return nil, err
	}

	if len(config.VaultToken) > 0 {
		err = v.TokenAuth(config.VaultToken)
	} else {
		err = v.AppIDAuth(config.VaultAppID, config.VaultUserIDPath)
	}

	return v, err
}
