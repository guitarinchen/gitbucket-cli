package org

import (
	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/org"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
	net_http "net/http"
)

var ListOrgsCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations for the authenticated user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return org.List(cmd, args, httpClient, presenterFactory)
	},
}
