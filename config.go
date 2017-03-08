package main

import (
	"fmt"
	"io"
)

// Config for polymerase
type Config struct {
	VaultAddr        string
	VaultToken       string
	VaultAppID       string
	VaultUserIDPath  string
	VaultFactoryFunc func(Config) (Vault, error)
	Input            io.Reader
	Output           io.Writer
}

// Validate the config
func (c Config) Validate() (bool, error) {
	if len(c.VaultAddr) == 0 {
		return false, fmt.Errorf("Invalid vault address")
	}

	if len(c.VaultToken) > 0 && (len(c.VaultAppID) > 0 || len(c.VaultUserIDPath) > 0) {
		return false, fmt.Errorf("Conflicting vault authentication strategies. Both app_id and token auth specified")
	}

	if len(c.VaultToken) == 0 && len(c.VaultAppID) == 0 && len(c.VaultUserIDPath) == 0 {
		return false, fmt.Errorf("No vault authentication strategy provided. Please specify a vault token or app ID and user ID path")
	}

	if (len(c.VaultAppID) > 0 && len(c.VaultUserIDPath) == 0) || (len(c.VaultAppID) == 0 && len(c.VaultUserIDPath) > 0) {
		return false, fmt.Errorf("Invalid vault authentication strategy provided. Please specify an app ID AND user ID path")
	}

	return true, nil
}
