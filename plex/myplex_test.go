package plex

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
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
	tokenFile = path.Join(".cache", "testTokenFile")
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
	tokenFile = path.Join(".cache", "testTokenFile")
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
	tokenFile = path.Join(".cache", "testTokenFile")
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
	tokenFile = path.Join(".cache", "testTokenFile")
	defer func() { tokenFile = oldtokenfile }()

	got := Token()
	if got != want {
		t.Errorf("Token() got %v want %v", got, want)
	}
}

// Test_apiRequest tests the apiRequest function
func Test_apiRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		token  string
		body   io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
		wantErr    bool
	}{
		{
			name: "Successful request",
			args: args{
				method: "GET",
				url:    "http://www.mocky.io/v2/5a07ad672f0000ae0ee610f1",
				token:  "testtoken",
				body:   nil,
			},
			wantStatus: http.StatusOK,
			wantBody:   "{ \"hello\": \"world\" }",
			wantErr:    false,
		},
		{
			name: "Bad URL",
			args: args{
				method: "GET",
				url:    "http://badurl",
				token:  "testtoken",
				body:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := apiRequest(tt.args.method, tt.args.url, tt.args.token, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.StatusCode != tt.wantStatus {
					t.Errorf("apiRequest() response status = %v, want %v", got.StatusCode, tt.wantStatus)
					return
				}

				defer got.Body.Close()
				body, err := ioutil.ReadAll(got.Body)
				if err != nil {
					t.Fatal("Error reading response body")
				}
				if string(body) != tt.wantBody {
					t.Errorf("apiRequest() body = %v, want %v", string(body), tt.wantBody)
				}
			}
		})
	}
}
