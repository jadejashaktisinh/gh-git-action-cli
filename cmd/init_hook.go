package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initHookCmd = &cobra.Command{
	Use:   "init-hook",
	Short: "Initialize a local git pre-push hook",
	Run: func(cmd *cobra.Command, args []string) {
		hookDir := filepath.Join(".git", "hooks")
		if _, err := os.Stat(hookDir); os.IsNotExist(err) {
			fmt.Println("Error: .git directory not found. Are you in a git repository?")
			return
		}

		hookPath := filepath.Join(hookDir, "pre-push")
		hookContent := `#!/bin/sh
# gh-git-action-cli pre-push hook
echo "🔍 Running local CI checks via gh-git-action-cli..."
gh-git-action-cli --job test
if [ $? -ne 0 ]; then
  echo "❌ Local CI checks failed. Push rejected."
  exit 1
fi
echo "✅ Local CI checks passed."
`

		err := os.WriteFile(hookPath, []byte(hookContent), 0755)
		if err != nil {
			fmt.Printf("Error writing hook: %v\n", err)
			return
		}

		fmt.Println("✅ Git pre-push hook initialized successfully!")
	},
}

func init() {
	rootCmd.AddCommand(initHookCmd)
}
