package cmd

import (
	"github.com/guitarinchen/gitbucket-cli/cmd/issue"
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issues",
	Long:  "Work with GitBucket issues.",
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.AddCommand(issue.CreateIssueCmd)
	issueCmd.AddCommand(issue.ListIssuesCmd)
}
