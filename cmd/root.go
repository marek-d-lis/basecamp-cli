package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "basecamp",
	Short: "BaseCamp CLI for setting up environments",
	Long:  "BaseCamp is a CLI tool for bootstrapping environments using Ansible.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
