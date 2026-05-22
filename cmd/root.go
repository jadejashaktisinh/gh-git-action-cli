/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jadejashaktisinh/gh-git-action-cli/config"
	"github.com/jadejashaktisinh/gh-git-action-cli/db"
	"github.com/jadejashaktisinh/gh-git-action-cli/parser"
	"github.com/jadejashaktisinh/gh-git-action-cli/runner"
	"github.com/jadejashaktisinh/gh-git-action-cli/tui"
	"github.com/spf13/cobra"
)

var (
	jobName string
	envFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-git-action-cli [workflow-file.yml]",
	Short: "Execute GitHub Actions workflows 100% locally",
	Long: `gh-git-action-cli is a GitHub CLI extension that intercepts and 
executes GitHub Actions workflows 100% locally on your computer.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && jobName == "" && envFile == "" {
			m := tui.InitialModel()
			p := tea.NewProgram(m, tea.WithAltScreen())
			finalModel, err := p.Run()
			if err != nil {
				fmt.Printf("Error running TUI: %v", err)
				os.Exit(1)
			}

			if tm, ok := finalModel.(tui.Model); ok && tm.Job != "" {
				jobName = tm.Job
				executeLocal([]string{tm.Workflow})
			}
			return
		}

		// Local execution
		executeLocal(args)
	},
}

func executeLocal(args []string) {
	workflowPath := ""
	if len(args) > 0 {
		workflowPath = args[0]
	} else {
		// Try to find workflows automatically
		files, _ := parser.FindWorkflows()
		if len(files) == 0 {
			fmt.Fprintln(os.Stderr, "Error: no workflow files found in .github/workflows/")
			os.Exit(1)
		}
		workflowPath = files[0] // Default to the first one for now
	}

	wf, err := parser.ParseWorkflow(workflowPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing workflow: %v\n", err)
		os.Exit(1)
	}

	env := make(map[string]string)
	if envFile != "" {
		loaded, err := parser.LoadEnvFile(envFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading env file: %v\n", err)
			os.Exit(1)
		}
		for k, v := range loaded {
			env[k] = v
		}
	}

	// Global workflow env
	for k, v := range wf.Env {
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	if jobName != "" {
		job, ok := wf.Jobs[jobName]
		if !ok {
			fmt.Fprintf(os.Stderr, "Error: job %s not found in %s\n", jobName, workflowPath)
			os.Exit(1)
		}

		// Merge job env
		for k, v := range job.Env {
			env[k] = v
		}

		err = runner.RunJob(job, runner.RunOptions{Env: env})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Job failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Run all jobs (sequentially for now)
		for _, job := range wf.Jobs {
			// Merge job env
			jobEnv := make(map[string]string)
			for k, v := range env {
				jobEnv[k] = v
			}
			for k, v := range job.Env {
				jobEnv[k] = v
			}

			err = runner.RunJob(job, runner.RunOptions{Env: jobEnv})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Job %s failed: %v\n", job.Name, err)
				os.Exit(1)
			}
		}
	}

	// Save to history
	record := db.Record{
		Timestamp:    time.Now(),
		WorkflowFile: workflowPath,
		TargetJob:    jobName,
		Conclusion:   "passed",
		Mode:         "native-shell",
		EnvSource:    envFile,
	}
	if err := db.SaveRun(record); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to save to history: %v\n", err)
	}
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
	rootCmd.Flags().StringVarP(&jobName, "job", "j", "", "Job name to execute")
	rootCmd.Flags().StringVarP(&envFile, "env-file", "e", "", "Path to .env file")
}
