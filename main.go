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
	Token    string
	Insecure bool
}

func getConfigFromEnv() Config {
	insecure, _ := strconv.ParseBool(os.Getenv("AKHQ_INSECURE"))
	return Config{
		BaseURL:  os.Getenv("AKHQ_URL"),
		Cluster:  os.Getenv("AKHQ_CLUSTER"),
		User:     os.Getenv("AKHQ_USER"),
		Pass:     os.Getenv("AKHQ_PASS"),
		Token:    os.Getenv("AKHQ_TOKEN"),
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

func readTopic(conf Config, client http.Client, topic string, token string) error {
	url := fmt.Sprintf("%s/api/%s/topic/%s/data", conf.BaseURL, conf.Cluster, topic)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v ", err)
	}
	req.Header.Set("Cookie", "JWT="+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error reading topic: %v", err)
	}
	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}
	fmt.Println(data.([]interface{}))

	//	for _, message := range data.([]interface{})["results"] {
	//		fmt.Println(message)
	//	}
	//
	// resp.Body.Close()
	// fmt.Println(data)
	return nil
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
			val := partition.(map[string]interface{})["offsetLag"]
			if val == nil {
				continue
			}
			if offsetLag, isFloat := val.(float64); isFloat {
				lag += uint(offsetLag)
			}
		}
		fmt.Println(lag)
		// fmt.Printf("%s\n\n", group["activeTopics"])
	}

}

func main() {
	conf := getConfigFromEnv()
	client := getClient(conf)
	token := conf.Token
	var err error
	if conf.Token == "" {
		token, err = Authenticate(conf, client)
		if err != nil {
			fmt.Println("Error authenticating: ", err)
			return
		}
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("AKHQ_TOKEN=%s", token)
		// fmt.Println("Topic is missing.")
		os.Exit(1)
	}
	GetOffsets(conf, client, args[0], token)
}
