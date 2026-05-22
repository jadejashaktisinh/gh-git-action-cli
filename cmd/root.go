/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/jadejashaktisinh/gh-git-action-cli/config"
	"github.com/jadejashaktisinh/gh-git-action-cli/db"
	"github.com/jadejashaktisinh/gh-git-action-cli/tui"
	"github.com/spf13/cobra"
)

var (
	repo    string
	branch  string
	inputs  []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-git-action-cli [workflow-file.yml]",
	Short: "Trigger, track, and manage manual GitHub Actions",
	Long: `gh-git-action-cli is a GitHub CLI extension to trigger, track, and manage 
manual GitHub Actions (workflow_dispatch) with local configuration, 
local history tracking, and an interactive terminal UI.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && repo == "" && branch == "" && len(inputs) == 0 {
			p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Error running TUI: %v", err)
				os.Exit(1)
			}
			return
		}

		// Headless execution
		executeHeadless(args)
	},
}

func executeHeadless(args []string) {
	if repo == "" {
		fmt.Fprintln(os.Stderr, "Error: --repo is required for headless execution")
		os.Exit(1)
	}

	workflow := ""
	if len(args) > 0 {
		workflow = args[0]
	} else {
		fmt.Fprintln(os.Stderr, "Error: workflow-file.yml is required as an argument")
		os.Exit(1)
	}

	if branch == "" {
		branch = "main" // Default to main
	}

	inputMap := make(map[string]string)
	for _, i := range inputs {
		// Expecting key=value
		parts := splitInput(i)
		if len(parts) == 2 {
			inputMap[parts[0]] = parts[1]
		}
	}

	fmt.Printf("Triggering workflow %s in %s on branch %s...\n", workflow, repo, branch)

	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating GitHub client: %v\n", err)
		os.Exit(1)
	}

	body := map[string]interface{}{
		"ref": branch,
	}
	if len(inputMap) > 0 {
		body["inputs"] = inputMap
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling request body: %v\n", err)
		os.Exit(1)
	}

	path := fmt.Sprintf("repos/%s/actions/workflows/%s/dispatches", repo, workflow)
	err = client.Post(path, bytes.NewReader(jsonBody), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error triggering workflow: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully triggered workflow!")

	// Save to history
	record := db.Record{
		Timestamp:  time.Now(),
		Repository: repo,
		Workflow:   workflow,
		Branch:     branch,
		Inputs:     inputMap,
		Conclusion: "triggered",
	}
	if err := db.SaveRun(record); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save to history: %v\n", err)
	}
}

func splitInput(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := config.InitConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.CloseDB()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "Repository to target (owner/repo)")
	rootCmd.PersistentFlags().StringVarP(&branch, "branch", "b", "", "Branch to target")
	rootCmd.Flags().StringSliceVarP(&inputs, "input", "i", []string{}, "Inputs in key=value format")
}


