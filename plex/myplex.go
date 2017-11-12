package plex

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"net/http"

	"encoding/xml"
)

// Credentials stores MyPlex credentials
type Credentials struct {
	Username string
	Password string
}

var (
	tokenFile = "token"
)

// Token requests a MyPlex authentication token from cache or from MyPlex.
func Token() string {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		// File does not exist. Get credentials and write token to file.
		log.Println("Cached token does not exist, request a new token in the Web Interface (Settings -> Request New Token).")
		return ""
	}
	log.Println("Using cached token.")
	return string(token)
}

// TokenRequest requests a new token from MyPlex
func TokenRequest(cred Credentials) error {
	type xmlUser struct {
		Email               string `xml:"email"`
		Username            string `xml:"username"`
		AuthenticationToken string `xml:"authentication-token"`
	}

	// Create a new reqest object.
	req, err := http.NewRequest("POST", "https://plex.tv/users/sign_in.xml", nil)
	if err != nil {
		return fmt.Errorf("failed to create new request")
	}
	// Configure the authentication and headers of the request.
	req.SetBasicAuth(cred.Username, cred.Password)
	addHeaders(*req, "")

	// Create the HTTP Client
	client := &http.Client{}

	// Get the response from the MyPlex API.
	log.Println("Requesting token from MyPlex servers.")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to MyPlex servers")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf(resp.Status)
	}

	var record xmlUser

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	err = xml.Unmarshal(body, &record)
	if err != nil {
		return fmt.Errorf("error parsing xml response: %v", err)
	}
	log.Println("Token received.")

	err = cacheToken(record.AuthenticationToken)
	if err != nil {
		return err
	}
	return nil
}

func cacheToken(token string) error {
	// Write token to file.
	f, err := os.Create(tokenFile)
	if err != nil {
		return fmt.Errorf("unable to create token file")
	}
	f.WriteString(token)
	f.Close()
	return nil
}

// RemoveCachedToken removes the cached tokenfile
func RemoveCachedToken() error {
	return os.Remove(tokenFile)
}

func addHeaders(r http.Request, token string) {
	r.Header.Add("X-Plex-Client-Identifier", "0bc797da-2ddd-4ce5-946e-5b13e48f17bb")
	r.Header.Add("X-Plex-Product", "Plex-Sync")
	r.Header.Add("X-Plex-Device", "Plex-Sync")
	r.Header.Add("X-Plex-Version", Version)
	r.Header.Add("X-Plex-Provides", "controller")
	if token != "" {
		r.Header.Add("X-Plex-Token", token)
	}
}

// GetToken requests the AccessToken from MyPlex for the named server
func (h *Host) GetToken(t string) error {
	// Create a new reqest object.
	resp, err := apiRequest("GET", "https://plex.tv/pms/servers.xml", t, nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}

	var record myPlexServer

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	err = xml.Unmarshal(body, &record)
	if err != nil {
		return fmt.Errorf("error parsing xml response: %v", err)
	}
	for _, x := range record.Server {
		if x.Name == h.Name {
			h.Token = x.AccessToken
			return nil
		}
	}
	return fmt.Errorf("no server found matching name %q", h.Name)
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

func apiRequest(method, url, token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request, %v", err)
	}
	addHeaders(*req, token)

	// Create the HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to MyPlex servers, %v", err)
	}
	return resp, nil
}
