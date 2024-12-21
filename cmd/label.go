package cmd

import (
	"github.com/guitarinchen/gitbucket-cli/cmd/label"
	"github.com/spf13/cobra"
)

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage labels",
	Long:  "Work with GitBucket labels.",
}

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(label.CreateLabelCmd)
	labelCmd.AddCommand(label.ListLabelsCmd)
}
