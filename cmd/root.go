package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type Config struct {
	BaseURL  string
	Cluster  string
	User     string
	Pass     string
	Token    string
	Insecure bool
}

var (
	conf   Config
	client http.Client
	token  string
)

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
	return http.Client{
		Transport: transp,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

var rootCmd = &cobra.Command{
	Use:   "akhq-cli",
	Short: "Query an AKHQ instance from the command line",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		conf = getConfigFromEnv()
		client = getClient(conf)
		token = conf.Token
		if token == "" {
			var err error
			token, err = Authenticate(conf, client)
			if err != nil {
				return fmt.Errorf("error authenticating: %v", err)
			}
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
