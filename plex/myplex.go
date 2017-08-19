package plex

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"net/http"

	"encoding/xml"

	"golang.org/x/crypto/ssh/terminal"
)

type credentials struct {
	username string
	password string
}

type prompter interface {
	promptCreds() credentials
}

type requester interface {
	tokenRequest(cred credentials) (string, error)
}

var (
	tokenFile = "token"
)

// Token requests a MyPlex authentication token from cache or from MyPlex.
func Token(pr prompter, r requester) (string, error) {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		// File does not exist. Get credentials and write token to file.
		log.Println("Cached token does not exist, prompt user for MyPlex credentials.")
		myplex := pr.promptCreds() // Get the user credentials.
		token, err := r.tokenRequest(myplex)
		if err != nil {
			return "", fmt.Errorf("error getting token: %v", err)
		}
		// Write token to file.
		f, err := os.Create(tokenFile)
		if err != nil {
			return "", fmt.Errorf("unable to create token file")
		}
		f.WriteString(token)
		f.Close()
		return token, nil
	}
	log.Println("Using cached token.")
	return string(token), nil
}

// CredPrompter is the method receiver for promptCreds
type CredPrompter struct{}

func (cp CredPrompter) promptCreds() credentials {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your MyPlex Username: ")
	user, _ := reader.ReadString('\n')

	fmt.Print("Enter your MyPlex Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	pass := string(bytePassword)
	fmt.Print("\n")

	return credentials{username: strings.TrimSpace(user), password: strings.TrimSpace(pass)}
}

// TokenRequester is the method receiver for tokenRequest
type TokenRequester struct{}

func (tr TokenRequester) tokenRequest(cred credentials) (string, error) {
	type xmlUser struct {
		Email               string `xml:"email"`
		Username            string `xml:"username"`
		AuthenticationToken string `xml:"authentication-token"`
	}

	// Create a new reqest object.
	req, err := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new request")
	}
	// Configure the authentication and headers of the request.
	req.SetBasicAuth(cred.username, cred.password)
	addHeaders(*req)

	// Create the HTTP Client
	client := &http.Client{}

	// Get the response from the MyPlex API.
	log.Println("Requesting token from MyPlex servers.")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed request to MyPlex servers")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf(string(http.StatusUnauthorized))
	}

	var record xmlUser

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	err = xml.Unmarshal(body, &record)
	if err != nil {
		return "", fmt.Errorf("error parsing xml response: %v", err)
	}
	log.Println("Token received.")

	return record.AuthenticationToken, nil
}

func addHeaders(r http.Request) {
	r.Header.Add("X-Plex-Client-Identifier", "0bc797da-2ddd-4ce5-946e-5b13e48f17bb")
	r.Header.Add("X-Plex-Product", "Plex-Sync")
	r.Header.Add("X-Plex-Device", "Plex-Sync")
	r.Header.Add("X-Plex-Version", Version)
	r.Header.Add("X-Plex-Provides", "controller")
}

// ServerAccessToken requests the AccessToken from MyPlex for the named server
func ServerAccessToken(t, name string) (string, error) {
	// Create a new reqest object.
	resp, err := apiRequest("GET", "https://plex.tv/pms/servers.xml?X-Plex-Token="+t, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	var record myPlexServer

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	err = xml.Unmarshal(body, &record)
	if err != nil {
		return "", fmt.Errorf("error parsing xml response: %v", err)
	}
	for _, x := range record.Server {
		if x.Name == name {
			return x.AccessToken, nil
		}
	}
	return "", fmt.Errorf("no server found matching name %q", name)
}

type myPlexServer struct {
	Server []struct {
		AccessToken    string `xml:"accessToken,attr"`
		Name           string `xml:"name,attr"`
		Address        string `xml:"address,attr"`
		Port           string `xml:"port,attr"`
		Version        string `xml:"version,attr"`
		Scheme         string `xml:"scheme,attr"`
		Host           string `xml:"host,attr"`
		LocalAddresses string `xml:"localAddresses,attr"`
		Owned          string `xml:"owned,attr"`
		Synced         string `xml:"synced,attr"`
	}
}

func apiRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request, %v", err)
	}
	addHeaders(*req)

	// Create the HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to MyPlex servers, %v", err)
	}
	return resp, nil
}

// CreateURI assembles the URI for an API request
func CreateURI(server Host, path, token string) string {
	if server.Ssl {
		return fmt.Sprintf("https://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
	}
	return fmt.Sprintf("http://%v:%v/%v?X-Plex-Token=%v", server.Hostname, server.Port, path, token)
}

// Host defines the data to be stored for server objects
type Host struct {
	Name     string
	Hostname string
	Port     int
	Ssl      bool
}
