package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tcurvelo/akhq-cli/internal/akhq"
)

// NewRootCmd builds the CLI; getenv is typically os.Getenv.
func NewRootCmd(args []string, getenv func(string) string, stdout, stderr io.Writer) *cobra.Command {
	client := &akhq.Client{} // configured in PersistentPreRunE, shared with subcommands

	root := &cobra.Command{
		Use:   "akhq-cli",
		Short: "Query an AKHQ instance from the command line",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			insecure, _ := strconv.ParseBool(getenv("AKHQ_INSECURE"))
			*client = *akhq.New(getenv("AKHQ_URL"), getenv("AKHQ_CLUSTER"), insecure)
			client.Token = getenv("AKHQ_TOKEN")
			if client.Token == "" {
				token, err := client.Login(getenv("AKHQ_USER"), getenv("AKHQ_PASS"))
				if err != nil {
					return fmt.Errorf("error authenticating: %v", err)
				}
				client.Token = token
			}
			return nil
		},
	}
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.AddCommand(newGroupsCmd(client), newSearchCmd(client))
	return root
}
