package repo

import (
	"fmt"
	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/repo"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
	net_http "net/http"
)

var ListRepoCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories owned by user or organization",
	Long:  "List repositories owned by a user or organization.",
	Args:  validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return repo.List(cmd, args, httpClient, presenterFactory)
	},
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("accepts at most 1 arg(s), received %d", len(args))
	}

	return nil
}
