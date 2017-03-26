package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type credentials struct {
	username string
	password string
}

type prompter interface {
	promptCreds() credentials
}

var (
	tokenFile = "token"
)

func token(pr prompter) string {
	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		// File does not exist. Get credentials and write token to file.
		log.Println("Cached token does not exist, prompt user for MyPlex credentials.")
		myplex := pr.promptCreds() // Get the user credentials.
		log.Printf("User: %s, Pass: %s", myplex.username, myplex.password)
		token := "token"
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

	return credentials{username: strings.TrimSpace(user), password: strings.TrimSpace(pass)}
}