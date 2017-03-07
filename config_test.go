package main

import "testing"

func TestValidateConfig(t *testing.T) {
	validWithToken := Config{VaultAddr: "google.com", VaultToken: "SomeToken"}
	validWithAppIDUserIDPath := Config{VaultAddr: "google.com", VaultAppID: "SomeID", VaultUserIDPath: "some/path"}
	invalidWithBoth := Config{VaultAddr: "google.com", VaultToken: "SomeToken", VaultAppID: "SomeID", VaultUserIDPath: "some/path"}
	invalidWithNeither := Config{VaultAddr: "google.com"}
	invalidWithNoAddr := Config{VaultAppID: "SomeID", VaultUserIDPath: "some/path"}
	invalidWithOnlyAppID := Config{VaultAddr: "google.com", VaultAppID: "SomeID"}
	invalidWithOnlyUserIDPath := Config{VaultAddr: "google.com", VaultUserIDPath: "some/path"}

	if valid, _ := validWithToken.Validate(); valid != true {
		t.Fatalf("Config %v was invalid but should have been valid", config)
	}

	if valid, _ := validWithAppIDUserIDPath.Validate(); valid != true {
		t.Fatalf("Config %v was invalid but should have been valid", config)
	}

	if valid, _ := invalidWithBoth.Validate(); valid != false {
		t.Fatalf("Config %v was valid but should have been invalid", config)
	}

	if valid, _ := invalidWithNeither.Validate(); valid != false {
		t.Fatalf("Config %v was valid but should have been invalid", config)
	}

	if valid, _ := invalidWithNoAddr.Validate(); valid != false {
		t.Fatalf("Config %v was valid but should have been invalid", config)
	}
	if valid, _ := invalidWithOnlyAppID.Validate(); valid != false {
		t.Fatalf("Config %v was valid but should have been invalid", config)
	}
	if valid, _ := invalidWithOnlyUserIDPath.Validate(); valid != false {
		t.Fatalf("Config %v was valid but should have been invalid", config)
	}
}
