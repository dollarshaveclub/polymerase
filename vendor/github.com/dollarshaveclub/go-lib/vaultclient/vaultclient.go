package vaultclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
)

const (
	authretries       = 10
	retrydelayseconds = 3
)

type VaultConfig struct {
	Server string // protocol, hostname and port (https://vault.foo.com:8200)
}

type VaultClient struct {
	client *api.Client
	config *VaultConfig
	token  string
}

// NewClient returns a VaultClient object or error
func NewClient(config *VaultConfig) (*VaultClient, error) {
	vc := VaultClient{}
	c, err := api.NewClient(&api.Config{Address: config.Server})
	vc.client = c
	vc.config = config
	return &vc, err
}

// TokenAuth sets the client token but doesn't check validity
func (c *VaultClient) TokenAuth(token string) error {
	c.token = token
	c.client.SetToken(token)
	ta := c.client.Auth().Token()
	var err error
	for i := 0; i < authretries; i++ {
		_, err = ta.LookupSelf()
		if err == nil {
			break
		}
		log.Printf("Token auth failed: %v, retrying (%v/%v)", err, i+1, authretries)
		time.Sleep(retrydelayseconds * time.Second)
	}
	if err != nil {
		return fmt.Errorf("error performing auth call to Vault (retries exceeded): %v", err)
	}
	return nil
}

// AppIDAuth attempts to perform app-id authorization.
func (c *VaultClient) AppIDAuth(appid string, useridpath string) error {
	f, err := os.Open(useridpath)
	if err != nil {
		return fmt.Errorf("error opening Vault User ID file: %v", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("error getting stat info on Vault User ID file: %v", err)
	}
	userid := make([]byte, fi.Size())
	_, err = f.Read(userid)
	if err != nil {
		return fmt.Errorf("error reading Vault User ID file: %v", err)
	}

	bodystruct := struct {
		AppID  string `json:"app_id"`
		UserID string `json:"user_id"`
	}{
		AppID:  appid,
		UserID: string(userid),
	}

	var resp *api.Response
	for i := 0; i < authretries; i++ {
		req := c.client.NewRequest("POST", "/v1/auth/app-id/login")
		jerr := req.SetJSONBody(bodystruct)
		if jerr != nil {
			return fmt.Errorf("error setting auth JSON body: %v", err)
		}
		resp, err = c.client.RawRequest(req)
		if err == nil {
			break
		}
		log.Printf("App-ID auth failed: %v, retrying (%v/%v)", err, i+1, authretries)
		time.Sleep(retrydelayseconds * time.Second)
	}
	if err != nil {
		return fmt.Errorf("error performing auth call to Vault (retries exceeded): %v", err)
	}

	var output interface{}
	jd := json.NewDecoder(resp.Body)
	err = jd.Decode(&output)
	if err != nil {
		return fmt.Errorf("error unmarshaling Vault auth response: %v", err)
	}
	body := output.(map[string]interface{})
	auth := body["auth"].(map[string]interface{})
	c.token = auth["client_token"].(string)
	return nil
}

// GetValue retrieves value at path
func (c *VaultClient) GetValue(path string) (interface{}, error) {
	c.client.SetToken(c.token)
	lc := c.client.Logical()
	s, err := lc.Read(path)
	if err != nil {
		return nil, fmt.Errorf("error reading secret from Vault: %v: %v", path, err)
	}
	if s == nil {
		return nil, fmt.Errorf("secret not found")
	}
	if _, ok := s.Data["value"]; !ok {
		return nil, fmt.Errorf("secret missing 'value' key")
	}
	return s.Data["value"], nil
}

// GetStringValue retrieves a value expected to be a string
func (c *VaultClient) GetStringValue(path string) (string, error) {
	val, err := c.GetValue(path)
	if err != nil {
		return "", err
	}
	switch val := val.(type) {
	case string:
		return val, nil
	default:
		return "", fmt.Errorf("unexpected type for %v value: %T", path, val)
	}
}

// GetBase64Value retrieves and decodes a value expected to be base64-encoded binary
func (c *VaultClient) GetBase64Value(path string) ([]byte, error) {
	val, err := c.GetStringValue(path)
	if err != nil {
		return []byte{}, err
	}
	decoded, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		return []byte{}, fmt.Errorf("vault path: %v: error decoding base64 value: %v", path, err)
	}
	return decoded, nil
}

// WriteValue writes value=data at path
func (c *VaultClient) WriteValue(path string, data []byte) error {
	c.client.SetToken(c.token)
	lc := c.client.Logical()
	_, err := lc.Write(path, map[string]interface{}{"value": data})
	return err
}
