package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcurvelo/akhq-cli/internal/akhq"
)

func newSearchCmd(client *akhq.Client) *cobra.Command {
	var params akhq.SearchParams
	searchCmd := &cobra.Command{
		Use:   "search <topic>",
		Short: "Search for data in a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := client.Search(args[0], params)
			if err != nil {
				return fmt.Errorf("error searching topic: %v", err)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), data)
			return err
		},
	}
	f := searchCmd.Flags()
	f.StringVar(&params.After, "after", "", "resume search after this key")
	f.Int32Var(&params.Partition, "partition", 0, "partition to search")
	f.StringVar(&params.Sort, "sort", "", "OLDEST or NEWEST")
	f.StringVar(&params.Timestamp, "timestamp", "", "search from this timestamp")
	f.StringVar(&params.SearchByKey, "search-by-key", "", "filter by message key")
	f.StringVar(&params.SearchByValue, "search-by-value", "", "filter by message value")
	f.StringVar(&params.SearchByHeaderKey, "search-by-header-key", "", "filter by header key")
	f.StringVar(&params.SearchByHeaderValue, "search-by-header-value", "", "filter by header value")
	return searchCmd
}
