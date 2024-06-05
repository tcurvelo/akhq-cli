package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	BaseURL  string
	Cluster  string
	User     string
	Pass     string
	Insecure bool
}

func getConfigFromEnv() Config {
	insecure, _ := strconv.ParseBool(os.Getenv("AKHQ_INSECURE"))
	return Config{
		BaseURL:  os.Getenv("AKHQ_URL"),
		Cluster:  os.Getenv("AKHQ_CLUSTER"),
		User:     os.Getenv("AKHQ_USER"),
		Pass:     os.Getenv("AKHQ_PASS"),
		Insecure: insecure,
	}
}

func getClient(conf Config) http.Client {
	transp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.Insecure},
	}
	client := &http.Client{
		Transport: transp,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return *client
}

func GetOffsets(conf Config, client http.Client, topic string, token string) {
	url := fmt.Sprintf("%s/api/%s/topic/%s/groups", conf.BaseURL, conf.Cluster, topic)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating new request:", err)
		return
	}
	req.Header.Set("Cookie", "JWT="+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error retrieving topic groups: ", err)
		return
	}
	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error reading topic: ", err)
		return
	}
	for _, group := range data.([]interface{}) {
		group := group.(map[string]interface{})
		// id := group["id"]
		var lag uint
		for _, partition := range group["offsets"].([]interface{}) {
			lag += uint(partition.(map[string]interface{})["offsetLag"].(float64))
		}
		fmt.Println(lag)
	}

}

func main() {
	conf := getConfigFromEnv()
	client := getClient(conf)
	token, err := Authenticate(conf, client)
	if err != nil {
		fmt.Println("Error authenticating: ", err)
		return
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Topic is missing.")
		os.Exit(1)
	}
	GetOffsets(conf, client, args[0], token)
}
