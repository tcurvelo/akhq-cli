package cmd

import (
	"bytes"
	"fmt"
	"net/http"
)

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
	_ = resp.Body.Close()

	return token, nil
}
