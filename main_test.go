package main

import (
	"os"
	"testing"
)

type fakeCredPrompter struct {
	username string
	password string
}

func (fcp fakeCredPrompter) promptCreds() credentials {
	return credentials{username: fcp.username, password: fcp.password}
}

func TestTokenFileRead(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	// Create a new temporary token file containing "testtoken".
	f, err := os.Create(tokenFile)
	if err != nil {
		t.Fatal("Unable to create token file.")
	}
	f.WriteString("testtoken")

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "TestUser",
		password: "TestPass",
	}

	// Check if the token function returns the value from the test token file.
	if token(fc) != "testtoken" {
		t.Error("Tokenfile does not contain 'testtoken'")
	}

	// Cleanup the temporary token file.
	f.Close()
	os.Remove(f.Name())
}

func TestTokenGeneration(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "TestUser",
		password: "TestPass",
	}

	// Check if the token function returns the value from the test token file.
	if token(fc) != "token" {
		t.Error("Generated token does not contain 'token'")
	}

	// Cleanup the temporary token file.
	f, err := os.Open(tokenFile)
	if err != nil {
		t.Fatal("Unable to remove test token file.")
	}
	f.Close()
	os.Remove(f.Name())
}
