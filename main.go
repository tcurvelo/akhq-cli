package main

import (
	"os"

	"github.com/tcurvelo/akhq-cli/cmd"
)

func main() {
	if err := cmd.NewRootCmd(os.Args[1:], os.Getenv, os.Stdout, os.Stderr).Execute(); err != nil {
		os.Exit(1)
	}
}
