package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "basecamp",
	Short: "⛺ BaseCamp CLI for setting up your coding camp",
	Long:  "BaseCamp CLI is your trusted tool for setting up and tearing down campsites for your development environments.",
}

func Execute() {
	printBanner() // Print the welcome message

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("🚫 Error encountered while setting up the camp:", err)
		os.Exit(1)
	}
}

func printBanner() {
	fmt.Println(`
  ⛺  Welcome to BaseCamp CLI  ⛺
---------------------------------
   "Prepare your camp and conquer!"`)
}
