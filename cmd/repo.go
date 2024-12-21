package cmd

import (
	"github.com/guitarinchen/gitbucket-cli/cmd/repo"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories",
	Long:  "Work with GitBucket repositories.",
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.AddCommand(repo.CreateRepoCmd)
	repoCmd.AddCommand(repo.ListRepoCmd)
	repoCmd.AddCommand(repo.CloneRepoCmd)
}
