package plex

import (
	"fmt"
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

func (ftr fakeTokenRequester) tokenRequest(cred credentials) (string, error) {
	if cred.username == "ValidUser" {
		return "ValidToken", nil
	} else if cred.username == "BadUser" {
		return "BadToken", fmt.Errorf("badtoken")
	} else {
		return "", fmt.Errorf("unknown")
	}
}

// Test the failure to create a token file.
func TestInvalidTokenFile(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "zzz:/invalidPath/tokenfile"
	defer func() { tokenFile = oldtokenfile }()

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "TestUser",
		password: "TestPass",
	}

	// Fake the token request.
	ft := fakeTokenRequester{}

	// Check if the token function returns the value from the test token file.
	_, err := Token(fc, ft)
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

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "TestUser",
		password: "TestPass",
	}

	// Fake the token request.
	ft := fakeTokenRequester{}

	// Check if the token function returns the value from the test token file.
	token, _ := Token(fc, ft)
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

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "ValidUser",
		password: "ValidPass",
	}

	// Fake the token request.
	ft := fakeTokenRequester{}

	// Check if the token function returns the value from the test token file.
	token, _ := Token(fc, ft)
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

// Test for invalid username/password.
func TestInvalidCredentials(t *testing.T) {
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	// Fake the credentials returned.
	fc := fakeCredPrompter{
		username: "Invalid",
		password: "Invalid",
	}

	tr := TokenRequester{}

	// Check if the token function returns the value from the test token file.
	_, err := Token(fc, tr)
	if err == nil {
		t.Error("Invalid credentials did not cause error")
	}
}

// TODO: Add test for blank token from token file.

// TODO: Add test for no token returned.

// TODO: Add test for no username/password supplied.