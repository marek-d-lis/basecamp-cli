package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Config holds the list of camp repositories from the YAML file
type Config struct {
	Camps []string `yaml:"camps"`
}

var (
	camps      []string
	configFile string
	tempCamp   = filepath.Join(os.TempDir(), "basecamp-campground")
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up your development camp",
	Long:  `"setup" sets up your development camp by setting up resources from the specified campsites.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load camps from config file if provided
		if configFile != "" {
			loadCampsFromFile(configFile)
		}

		// Error if no camps were provided
		if len(camps) == 0 {
			fmt.Println("🚫 No campsites specified. Please provide at least one campsite using --camp or --config.")
			fmt.Println("Use './basecamp setup --help' for more information.")
			os.Exit(1)
		}

		// Ensure campground exists
		if err := os.MkdirAll(tempCamp, 0755); err != nil {
			fmt.Printf("🚫 Could not prepare the campground: %v\n", err)
			os.Exit(1)
		}

		// Pitch the camp for each campsite
		fmt.Printf("⛺ Setting up camp from %d campsites...\n", len(camps))
		for _, camp := range camps {
			setupCamp(camp)
		}

		// Clean up the campground after setting up
		breakCamp()

		fmt.Println("🔥 The camp is ready! Your environment is fully prepared.")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.Flags().StringSliceVarP(&camps, "camp", "c", []string{}, "Specify campsites to set up (can specify multiple)")
	setupCmd.Flags().StringVarP(&configFile, "config", "f", "", "Specify a YAML config file containing campsites")
}

func loadCampsFromFile(file string) {
	fmt.Printf("🗺️  Loading campsite map from %s...\n", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("🚫 Could not read campsite map: %v\n", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("🚫 Could not understand the campsite map: %v\n", err)
		os.Exit(1)
	}

	camps = append(camps, config.Camps...)
}

func setupCamp(camp string) {
	campName := filepath.Base(camp)
	campPath := filepath.Join(tempCamp, campName)

	fmt.Printf("⛏️  Preparing camp: %s...\n", campName)
	if _, err := os.Stat(campPath); os.IsNotExist(err) {
		fmt.Printf("🔄 Gathering supplies from %s...\n", camp)
		if err := exec.Command("git", "clone", camp, campPath).Run(); err != nil {
			fmt.Printf("🚫 Failed to gather supplies from %s.\n", camp)
			os.Exit(1)
		}
	} else {
		fmt.Printf("🔄 Updating existing camp: %s...\n", campName)
		if err := exec.Command("git", "-C", campPath, "pull").Run(); err != nil {
			fmt.Printf("🚫 Could not update camp: %s.\n", campName)
			os.Exit(1)
		}
	}

	playbookPath := filepath.Join(campPath, "camp.yml")
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  No camp setup instructions found in %s. Skipping...\n", campName)
		return
	}

	fmt.Printf("🔥 Building the camp from %s...\n", campName)
	cmd := exec.Command("ansible-playbook", playbookPath, "--ask-become-pass")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("🚫 Could not set up the camp: %v\n", err)
		os.Exit(1)
	}
}

func breakCamp() {
	fmt.Printf("🧹 Breaking camp and cleaning up: %s\n", tempCamp)
	if err := os.RemoveAll(tempCamp); err != nil {
		fmt.Printf("🚫 Could not clean up the campground: %v\n", err)
		os.Exit(1)
	}
}
