package issue

import (
	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/issue"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
	net_http "net/http"
)

var CreateIssueCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  "Create an issue on GitBucket.",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return issue.Create(cmd, args, httpClient, presenterFactory)
	},
}

func init() {
	CreateIssueCmd.Flags().StringVarP(&issue.CreateFlags.Repo, "repo", "r", issue.CreateFlags.Repo, "Select a repository using the [OWNER/]REPO format")
	CreateIssueCmd.Flags().StringVarP(&issue.CreateFlags.Title, "title", "t", issue.CreateFlags.Title, "Supply a title. Will prompt for one otherwise.")
	CreateIssueCmd.Flags().StringVarP(&issue.CreateFlags.Body, "body", "b", issue.CreateFlags.Body, "Supply a body. Will prompt for one otherwise.")
	CreateIssueCmd.Flags().StringArrayVarP(&issue.CreateFlags.Assignees, "assignees", "a", issue.CreateFlags.Assignees, "Assign people by their login. Use \"@me\" to self-assign.")
	CreateIssueCmd.Flags().IntVarP(&issue.CreateFlags.Milestone, "milestone", "m", issue.CreateFlags.Milestone, "Add the issue to a milestone by name")
	CreateIssueCmd.Flags().StringArrayVarP(&issue.CreateFlags.Labels, "labels", "l", issue.CreateFlags.Labels, "Add labels by name")
}
