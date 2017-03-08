package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	template := "{{ .FIRST_NAME }}"
	output := &bytes.Buffer{}
	context := newTestContext("", template, output)
	setupTest(context)

	os.Setenv("FIRST_NAME", "JAMES")
	run(rootCmd, []string{})
	validateOutput(output, "JAMES", t)
}

func TestVault(t *testing.T) {
	template := "{{ vault \"secret_agents/007/last_name\" }}"
	output := &bytes.Buffer{}
	context := newTestContext("BOND", template, output)
	setupTest(context)

	run(rootCmd, []string{})
	validateOutput(output, "BOND", t)
}

func TestFile(t *testing.T) {
	output := &bytes.Buffer{}
	context := newTestContext("BOND", "", output)
	setupTest(context)

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile(wd, "polymerase_test_template")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	_, _ = tmpfile.WriteString("{{ .FIRST_NAME }} {{ vault \"secret_agents/007/last_name\" }}")
	filename := tmpfile.Name()
	_ = tmpfile.Close()

	run(rootCmd, []string{filename})
	validateOutput(output, "JAMES BOND", t)
}

func validateOutput(actual *bytes.Buffer, expected string, t *testing.T) {
	outStr := string(actual.Bytes())
	if outStr != expected {
		t.Fatalf("Expected %v but got %v", expected, outStr)
	}
}

func setupTest(context *testContext) {
	config = newTestConfig(context.mockVault.Vault, context.inputBuffer, context.outputBuffer)
}

type testContext struct {
	mockVault    *mockVaultClient
	inputBuffer  *bytes.Buffer
	outputBuffer *bytes.Buffer
}

func newTestContext(vaultValue string, template string, outputBuffer *bytes.Buffer) *testContext {
	return &testContext{mockVault: &mockVaultClient{value: vaultValue},
		inputBuffer:  bytes.NewBufferString(template),
		outputBuffer: outputBuffer,
	}
}

func newTestConfig(vf func(Config) (Vault, error), input io.Reader, output io.Writer) Config {
	return Config{VaultAddr: "ADDR", VaultToken: "TESTTOKEN", VaultFactoryFunc: vf, Input: input, Output: output}
}

type mockVaultClient struct {
	value string
}

func (c mockVaultClient) GetStringValue(path string) (string, error) {
	return c.value, nil
}

func (c mockVaultClient) Vault(config Config) (Vault, error) {
	return c, nil
}
