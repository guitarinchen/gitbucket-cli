package issue

import (
	net_http "net/http"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/issue"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

var ListIssuesCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues in a repository",
	Long:  "List issues in a GitBucket repository.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return issue.List(cmd, args, httpClient, presenterFactory)
	},
}

func init() {
	ListIssuesCmd.Flags().StringVarP(&issue.ListFlags.State, "state", "s", "open", "Filter by state: {open|closed|all} (default \"open\")")
}
