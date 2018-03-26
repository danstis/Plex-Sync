package plex

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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

// Test removing a new token file.
func TestTokenRemoval(t *testing.T) {
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
	f.Close()

	if err := RemoveCachedToken(); err != nil {
		t.Error("Error removing token file")
	}
}

// Test missing token file
func TestNewTokenfile(t *testing.T) {
	want := ""
	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() { tokenFile = oldtokenfile }()

	got := Token()
	if got != want {
		t.Errorf("Token() got %v want %v", got, want)
	}
}

func TestMyPlexToken(t *testing.T) {

	ts := startMyPlexTestServer()
	defer ts.Close()

	// Replace the tokenfile path for the duration of this test.
	oldtokenfile := tokenFile
	tokenFile = "testTokenFile"
	defer func() {
		os.Remove(tokenFile)
		tokenFile = oldtokenfile
	}()

	if tokenRequest(Credentials{}, ts.URL+"/users/goodSign_in.xml") != nil {
		t.Errorf("Testing a good signin failed")
	}
	if tokenRequest(Credentials{}, ts.URL+"/users/badSign_in.xml") == nil {
		t.Errorf("Testing a bad signin did not result in an Error")
	}
}

func startMyPlexTestServer() *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {

		case "/users/goodSign_in.xml":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <user>
                    <username>testuser</username>
                    <email>testuser@domain.com</email>
                    <joined-at type="datetime">2010-12-31 12:59:59 UTC</joined-at>
                    <authentication-token>GoodToken</authentication-token>
                </user>`)

		case "/users/badSign_in.xml":
			w.Header().Set("Content-Type", "text/xml;charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
                <errors>
                    <error>Invalid email, username, or password.</error>
                </errors>`)

		}

	}))

	return s
}
