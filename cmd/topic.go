package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	return nil
}
