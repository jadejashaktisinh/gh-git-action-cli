package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <run-id>",
	Short: "Query real-time status of a workflow run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runID := args[0]
		
		targetRepo := repo
		if targetRepo == "" {
			fmt.Fprintln(os.Stderr, "Error: repository is required. Use --repo flag.")
			return
		}

		client, err := api.DefaultRESTClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating GitHub client: %v\n", err)
			return
		}

		var run struct {
			Status     string `json:"status"`
			Conclusion string `json:"conclusion"`
		}

		path := fmt.Sprintf("repos/%s/actions/runs/%s", targetRepo, runID)
		err = client.Get(path, &run)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching run status: %v\n", err)
			return
		}

		fmt.Printf("Run ID: %s\n", runID)
		fmt.Printf("Status: %s\n", run.Status)
		if run.Conclusion != "" {
			fmt.Printf("Conclusion: %s\n", run.Conclusion)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
