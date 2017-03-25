package main

import (
	"os"
	"testing"
)

func TestToken(t *testing.T) {
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

	// Check if the token function returns the value from the test token file.
	if token() != "testtoken" {
		t.Error("Tokenfile does not contain 'testtoken'")
	}

	// Cleanup the temporary token file.
	f.Close()
	os.Remove(f.Name())
}
