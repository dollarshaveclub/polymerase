package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dollarshaveclub/polymerase/pkg/vaultclient"
	"github.com/spf13/cobra"
)

var vault *vaultclient.VaultClient
var logger = log.New(os.Stderr, "", log.LstdFlags)

// Config for polymerase
type Config struct {
	VaultAddr       string
	VaultToken      string
	VaultAppID      string
	VaultUserIDPath string
}

var config = Config{}

var rootCmd = &cobra.Command{
	Use:     "polymerase",
	Short:   "polymerase",
	Long:    "Templates a file at the specified path using environment variables and vault values.",
	Example: "polymerase <filename>",
	Run:     run,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.VaultAppID, "app-id", "a", os.Getenv("APP_ID"), "Vault App-ID. Can use APP_ID environment variable instead.")
	rootCmd.PersistentFlags().StringVarP(&config.VaultAddr, "vault-addr", "v", os.Getenv("VAULT_ADDR"), "Vault server address (including protocol and port). Can use VAULT_ADDR environment variable instead.")
	rootCmd.PersistentFlags().StringVarP(&config.VaultToken, "vault-token", "t", os.Getenv("VAULT_TOKEN"), "Vault token. Can use VAULT_TOKEN environment variable instead.")
	rootCmd.PersistentFlags().StringVarP(&config.VaultUserIDPath, "user-id-path", "u", os.Getenv("USER_ID_PATH"), "Path to user id. Can use USER_ID_PATH environment variable instead.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {

	if len(args) != 1 {
		cmd.Usage()
		return
	}

	if err := validateConfig(config); err != nil {
		logger.Fatalf("Error validating config: %v", err)
	}

	var err error
	vault, err = authenticatedVaultUsingConfig(config)
	if err != nil {
		logger.Fatalf("Error configuring vault: %v", err)
	}

	funcMap := template.FuncMap{"vault": vaultGetString}
	filename := args[0]
	tplName := filepath.Base(filename)
	tmpl, err := template.New(tplName).Funcs(funcMap).ParseFiles(filename)
	if err != nil {
		logger.Fatalf("Error parsing template: %v", err)
	}

	err = tmpl.Execute(os.Stdout, env())
	if err != nil {
		logger.Fatalf("Error populating template: %v", err)
	}
}

func validateConfig(config Config) error {
	if len(config.VaultAddr) == 0 {
		return fmt.Errorf("Invalid vault address")
	}

	if len(config.VaultToken) > 0 && (len(config.VaultAppID) > 0 || len(config.VaultUserIDPath) > 0) {
		return fmt.Errorf("Conflicting vault authentication strategies. Both app_id and token auth specified")
	}

	if len(config.VaultToken) == 0 && len(config.VaultAppID) == 0 && len(config.VaultUserIDPath) == 0 {
		return fmt.Errorf("No vault authentication strategy provided. Please specify a vault token or app ID and user ID")
	}

	return nil
}

func authenticatedVaultUsingConfig(config Config) (*vaultclient.VaultClient, error) {

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

func env() map[string]string {
	env := make(map[string]string)
	for _, item := range os.Environ() {
		key, val := envKeyVal(item)
		env[key] = val
	}

	return env
}

func envKeyVal(env string) (string, string) {
	spl := strings.Split(env, "=")

	return spl[0], strings.Join(spl[1:], "=")
}

func vaultGetString(path string) string {
	val, err := vault.GetStringValue(path)
	if err != nil {
		logger.Fatalf("Error fetching value from vault: %v", err)
	}

	return val
}
