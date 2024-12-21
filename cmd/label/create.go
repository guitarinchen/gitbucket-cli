package label

import (
	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/label"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
	net_http "net/http"
)

var CreateLabelCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	Long:  "Create a label on GitBucket.",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return label.Create(cmd, args, httpClient, presenterFactory)
	},
}

func init() {
	CreateLabelCmd.Flags().StringVarP(&label.CreateFlags.Repo, "repo", "r", "", "Select a repository using the [OWNER/]REPO format")
	CreateLabelCmd.Flags().StringVarP(&label.CreateFlags.Name, "name", "n", "", "Color of the label")
	CreateLabelCmd.Flags().StringVarP(&label.CreateFlags.Color, "color", "c", "", "Color of the label")
}
