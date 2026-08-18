package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

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
	}
}

var groupsCmd = &cobra.Command{
	Use:   "groups <topic>",
	Short: "Print total consumer group lag for a topic",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		GetOffsets(conf, client, args[0], token)
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
