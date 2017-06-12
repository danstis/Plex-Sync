package plex

import (
	"log"
	"os"
	"testing"
)

type fakeCredPrompter struct {
	username string
	password string
}

type fakeTokenRequester struct {
}

func (fcp fakeCredPrompter) promptCreds() credentials {
	return credentials{username: fcp.username, password: fcp.password}
}

func (ftr fakeTokenRequester) tokenRequest(cred credentials) string {
	if cred.username == "ValidUser" {
		return "ValidToken"
	} else if cred.username == "BadUser" {
		return "BadToken"
	} else {
		return "token"
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

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "TestUser",
		password: "TestPass",
	}

	// Fake the token request.
	ft := fakeTokenRequester{}

	// Check if the token function returns the value from the test token file.
	if token(fc, ft) != "ValidToken" {
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

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "ValidUser",
		password: "ValidPass",
	}

	// Fake the token request.
	ft := fakeTokenRequester{}

	// Check if the token function returns the value from the test token file.
	if token(fc, ft) != "ValidToken" {
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

// TODO: Add test for invalid username/password.

// TODO: Add test for blank token from token file.

// TODO: Add test for no token returned.

// TODO: Add test for no username/password supplied.
