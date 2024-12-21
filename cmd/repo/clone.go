package repo

import (
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/repo"
	"github.com/spf13/cobra"
)

var CloneRepoCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone/sync with a remote repository",
	Long:  "Clone a GitBucket repository locally.",
	Args:  cobra.ExactArgs(1),
	RunE:  repo.Clone,
}

func init() {
	CloneRepoCmd.Flags().BoolVarP(&repo.CloneFlags.UseSSH, "ssh", "s", false, "Clone with SSH (default: false)")
	CloneRepoCmd.Flags().StringVarP(&repo.CloneFlags.Branch, "branch", "b", "master", "Specify branch name.")
}
