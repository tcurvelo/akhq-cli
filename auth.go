package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func Authenticate(conf Config, client http.Client) (string, error) {
	var token string
	url := fmt.Sprintf("%s/login", conf.BaseURL)
	body := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, conf.User, conf.Pass)
	payload := bytes.NewBuffer([]byte(body))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "JWT" {
			token = cookie.Value
			break
		}
	}
	resp.Body.Close()

	return token, nil
}
