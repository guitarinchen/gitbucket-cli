package cmd

import (
	"github.com/guitarinchen/gitbucket-cli/cmd/org"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage organizations",
	Long:  "Work with GitBucket organizations.",
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(org.ListOrgsCmd)
	orgCmd.AddCommand(org.CreateOrgCmd)
}
