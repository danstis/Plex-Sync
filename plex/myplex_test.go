package plex

import (
	"log"
	"os"
	"testing"
)

// Test the failure to create a token file.
func TestInvalidTokenFile(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "zzz:/invalidPath/tokenfile"
	defer func() { tokenFile = oldtokenfile }()

	// Check if the token function returns the value from the test token file.
	err := cacheToken("test")
	if err == nil {
		t.Error("Was able to create invalid token file")
	}
}

// Test reading of the token from a temporary token file.
func TestTokenFileRead(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	// Create a new temporary token file containing "ValidToken".
	f, err := os.Create(tokenFile)
	if err != nil {
		t.Fatal("Unable to create token file.")
	}
	f.WriteString("ValidToken")

	// Check if the token function returns the value from the test token file.
	token := Token()
	if token != "ValidToken" {
		t.Error("Tokenfile does not contain 'ValidToken'")
	}

	// Cleanup the temporary token file.
	f.Close()
	os.Remove(f.Name())
}

// Test getting a new token with a valid username and password, then writing it to file.
func TestTokenGeneration(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	err := cacheToken("ValidToken")
	if err != nil {
		t.Error("Unable to generate tokenfile")
	}

	// Check if the token function returns the value from the test token file.
	token := Token()
	if token != "ValidToken" {
		t.Error("Generated token does not contain 'ValidToken'")
	}

	// Cleanup the temporary token file.
	f, err := os.Open(tokenFile)
	if err != nil {
		log.Println(err)
	}
	f.Close()
	if err := os.Remove(f.Name()); err != nil {
		log.Printf("Error removing file: %s", err)
	}
}
