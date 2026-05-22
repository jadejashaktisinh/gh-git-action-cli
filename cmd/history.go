package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jadejashaktisinh/gh-git-action-cli/db"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show local history of workflow runs",
	Run: func(cmd *cobra.Command, args []string) {
		records, err := db.GetHistory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching history: %v\n", err)
			return
		}

		if len(records) == 0 {
			fmt.Println("No history found.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIMESTAMP\tWORKFLOW\tJOB\tMODE\tSTATUS")
		
		for _, r := range records {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				r.Timestamp.Format(time.RFC3339),
				r.WorkflowFile,
				r.TargetJob,
				r.Mode,
				r.Conclusion,
			)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
