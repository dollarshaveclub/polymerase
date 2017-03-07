package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var vault Vault
var logger = log.New(os.Stderr, "", log.LstdFlags)
var config = Config{VaultFactoryFunc: AuthenticatedVaultClient, Input: os.Stdin, Output: os.Stdout}

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

	if len(args) > 1 {
		cmd.Usage()
		return
	}

	if _, err := config.Validate(); err != nil {
		logger.Fatalf("Error validating config: %v", err)
	}

	var err error
	vault, err = config.VaultFactoryFunc(config)
	if err != nil {
		logger.Fatalf("Error configuring vault: %v", err)
	}

	var tmpl Template
	if len(args) == 1 {
		tmpl, err = TemplateFromFile(args[0])
	} else {
		tmpl, err = TemplateFromReader(config.Input)
	}

	if err != nil {
		logger.Fatalf("Error parsing template: %v", err)
	}

	err = tmpl.Execute(config.Output, env())
	if err != nil {
		logger.Fatalf("Error populating template: %v", err)
	}
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
