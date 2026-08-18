package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

type SearchParams struct {
	After               string
	Partition           int32
	Sort                string // OLDEST or NEWEST
	Timestamp           string
	SearchByKey         string
	SearchByValue       string
	SearchByHeaderKey   string
	SearchByHeaderValue string
}

func searchTopic(conf Config, client http.Client, topic string, token string, params SearchParams) error {
	base := fmt.Sprintf("%s/api/%s/topic/%s/data/search", conf.BaseURL, conf.Cluster, topic)

	q := url.Values{}
	q.Set("partition", strconv.Itoa(int(params.Partition)))
	for key, val := range map[string]string{
		"after":               params.After,
		"sort":                params.Sort,
		"timestamp":           params.Timestamp,
		"searchByKey":         params.SearchByKey,
		"searchByValue":       params.SearchByValue,
		"searchByHeaderKey":   params.SearchByHeaderKey,
		"searchByHeaderValue": params.SearchByHeaderValue,
	} {
		if val != "" {
			q.Set(key, val)
		}
	}

	req, err := http.NewRequest("GET", base+"?"+q.Encode(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v ", err)
	}
	req.Header.Set("Cookie", "JWT="+token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error searching topic: %v", err)
	}
	var data interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error parsing response: %v", err)
	}
	fmt.Println(data)
	return nil
}

var searchParams SearchParams

var searchCmd = &cobra.Command{
	Use:   "search <topic>",
	Short: "Search for data in a topic",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return searchTopic(conf, client, args[0], token, searchParams)
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchParams.After, "after", "", "resume search after this key")
	searchCmd.Flags().Int32Var(&searchParams.Partition, "partition", 0, "partition to search")
	searchCmd.Flags().StringVar(&searchParams.Sort, "sort", "", "OLDEST or NEWEST")
	searchCmd.Flags().StringVar(&searchParams.Timestamp, "timestamp", "", "search from this timestamp")
	searchCmd.Flags().StringVar(&searchParams.SearchByKey, "search-by-key", "", "filter by message key")
	searchCmd.Flags().StringVar(&searchParams.SearchByValue, "search-by-value", "", "filter by message value")
	searchCmd.Flags().StringVar(&searchParams.SearchByHeaderKey, "search-by-header-key", "", "filter by header key")
	searchCmd.Flags().StringVar(&searchParams.SearchByHeaderValue, "search-by-header-value", "", "filter by header value")
	rootCmd.AddCommand(searchCmd)
}
