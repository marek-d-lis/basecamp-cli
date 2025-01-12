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

// Config holds the repositories from the YAML file
type Config struct {
	Repos []string `yaml:"repos"`
}

var (
	repos      []string
	configFile string
	tempDir    = filepath.Join(os.TempDir(), "basecamp-cli")
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the BaseCamp setup process",
	Long: `The "run" command runs the specified repositories' Ansible playbooks.
You must specify repositories via --repo flags or use a basecamp.yml configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load repositories from config file if provided
		if configFile != "" {
			loadConfigFromFile(configFile)
		}

		// Error if no repositories were provided
		if len(repos) == 0 {
			fmt.Println("❌ No repositories specified. Please provide at least one repository using --repo or --config.")
			fmt.Println("Use './basecamp run --help' for more information.")
			os.Exit(1)
		}

		// Ensure temp directory exists
		if err := os.MkdirAll(tempDir, 0755); err != nil {
			fmt.Printf("❌ Failed to create temp directory: %v\n", err)
			os.Exit(1)
		}

		// Process each repository
		fmt.Printf("📦 Cloning and running playbooks for %d repositories...\n", len(repos))
		for _, repo := range repos {
			processRepository(repo)
		}

		// Clean up temporary directory after running
		cleanupTempDir()

		fmt.Println("🚀 All playbooks completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringSliceVarP(&repos, "repo", "r", []string{}, "Add repositories to run (can specify multiple)")
	runCmd.Flags().StringVarP(&configFile, "config", "c", "", "Specify a YAML config file containing repositories")
}

func loadConfigFromFile(file string) {
	fmt.Printf("📖 Loading configuration from %s...\n", file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("❌ Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("❌ Failed to parse YAML: %v\n", err)
		os.Exit(1)
	}

	repos = append(repos, config.Repos...)
}

func processRepository(repo string) {
	repoName := filepath.Base(repo)
	repoPath := filepath.Join(tempDir, repoName)

	fmt.Printf("📦 Processing repository: %s\n", repo)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Cloning %s into %s...\n", repo, repoPath)
		if err := exec.Command("git", "clone", repo, repoPath).Run(); err != nil {
			fmt.Printf("❌ Failed to clone repository: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("✅ Repository already exists. Pulling latest changes...\n")
		if err := exec.Command("git", "-C", repoPath, "pull").Run(); err != nil {
			fmt.Printf("❌ Failed to pull changes: %v\n", err)
			os.Exit(1)
		}
	}

	playbookPath := filepath.Join(repoPath, "dev_setup.yml")
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  No playbook found in %s. Skipping...\n", repo)
		return
	}

	fmt.Printf("🚀 Running playbook: %s\n", playbookPath)
	cmd := exec.Command("ansible-playbook", playbookPath, "--ask-become-pass")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Failed to run playbook: %v\n", err)
		os.Exit(1)
	}
}

func cleanupTempDir() {
	fmt.Printf("🧹 Cleaning up temporary directory: %s\n", tempDir)
	if err := os.RemoveAll(tempDir); err != nil {
		fmt.Printf("❌ Failed to clean up temp directory: %v\n", err)
		os.Exit(1)
	}
}
