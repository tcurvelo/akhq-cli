package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tcurvelo/akhq-cli/internal/akhq"
)

func newGroupsCmd(client *akhq.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "groups <topic>",
		Short: "Print total consumer group lag for a topic",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			groups, err := client.TopicGroups(args[0])
			if err != nil {
				// errors on stdout with exit 0 are kept until the output contract (#6)
				_, _ = fmt.Fprintln(out, "Error retrieving topic groups: ", err)
				return
			}
			for _, group := range groups {
				var offsetLag int64
				for _, p := range group.Offsets {
					if p.OffsetLag != nil {
						offsetLag += *p.OffsetLag
					}
				}
				_, _ = fmt.Fprintln(out, offsetLag)
			}
		},
	}
}
