package main

import (
	"bufio"
	"fmt"
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
	tokenRequest(cred credentials) string
}

var (
	tokenFile = "token"
)

func token(pr prompter, r requester) string {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		// File does not exist. Get credentials and write token to file.
		log.Println("Cached token does not exist, prompt user for MyPlex credentials.")
		myplex := pr.promptCreds() // Get the user credentials.
		token := r.tokenRequest(myplex)
		// Write token to file.
		f, err := os.Create(tokenFile)
		if err != nil {
			log.Fatalln("Error: Unable to create token file.")
		}
		f.WriteString(token)
		f.Close()
		return token
	}
	log.Println("Using cached token.")
	return string(token)
}

type credPrompter struct{}

func (cp credPrompter) promptCreds() credentials {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your MyPlex Username: ")
	user, _ := reader.ReadString('\n')

	fmt.Print("Enter your MyPlex Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	pass := string(bytePassword)
	fmt.Print("\n")

	return credentials{username: strings.TrimSpace(user), password: strings.TrimSpace(pass)}
}

type tokenRequester struct{}

func (tr tokenRequester) tokenRequest(cred credentials) string {
	type XMLUser struct {
		Email               string `xml:"email"`
		Username            string `xml:"username"`
		AuthenticationToken string `xml:"authentication-token"`
	}

	// The MyPlex API URL for performing signin.
	uri := "https://plex.tv/users/sign_in.xml"

	// Create a new reqest object.
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	// Configure the authentication and headers of the request.
	req.SetBasicAuth(cred.username, cred.password)
	req.Header.Add("X-Plex-Client-Identifier", "0bc797da-2ddd-4ce5-946e-5b13e48f17bb")
	req.Header.Add("X-Plex-Product", "Plex-Sync")
	req.Header.Add("X-Plex-Device", "Plex-Sync")
	req.Header.Add("X-Plex-Version", Version)
	req.Header.Add("X-Plex-Provides", "controller")

	// Create the HTTP Client
	client := &http.Client{}

	// Get the response from the MyPlex API.
	log.Println("Requesting token from MyPlex servers.")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close() // Close the connection once reading is complete.

	var record XMLUser

	if body, readerr := ioutil.ReadAll(resp.Body); readerr == nil {
		//load object xml
		if xmlerr := xml.Unmarshal(body, &record); xmlerr != nil {
			log.Println(xmlerr)
		}
	} else {
		//return read error
		log.Println(readerr)
	}
	log.Println("Token received.")

	return record.AuthenticationToken
}
