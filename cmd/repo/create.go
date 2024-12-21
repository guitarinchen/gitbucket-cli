package repo

import (
	net_http "net/http"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/repo"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

var CreateRepoCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repository",
	Long:  "Create a new GitBucket repository.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return repo.Create(cmd, args, httpClient, presenterFactory)
	},
}

func init() {
	CreateRepoCmd.Flags().BoolVarP(&repo.CreateFlags.IsPrivate, "private", "p", false, "Make the new repository private")
	CreateRepoCmd.Flags().StringVarP(&repo.CreateFlags.Description, "description", "d", "", "Description of the repository")
	CreateRepoCmd.Flags().BoolVarP(&repo.CreateFlags.AutoInit, "auto-init", "a", false, "Add a README file to the new repository")
}
