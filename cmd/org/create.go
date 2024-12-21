package org

import (
	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/org"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
	net_http "net/http"
)

var CreateOrgCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an organization.",
	RunE: func(cmd *cobra.Command, args []string) error {
		httpClient := http.NewClient(configs.ApiToken)
		presenterFactory := func(resp *net_http.Response) presenter.Presenter {
			return presenter.NewResponsePresenter(resp)
		}
		return org.Create(cmd, args, httpClient, presenterFactory)
	},
}

func init() {
	CreateOrgCmd.Flags().StringVarP(&org.CreateFlags.Login, "login", "l", "", "The organization's username")
	CreateOrgCmd.Flags().StringVarP(&org.CreateFlags.Admin, "admin", "a", "", "The login of the user who will manage this organization")
	CreateOrgCmd.Flags().StringVarP(&org.CreateFlags.ProfileName, "profile-name", "p", "", "The organization's display name")
	CreateOrgCmd.Flags().StringVarP(&org.CreateFlags.URL, "url", "u", "", "The organization's avator url")
}
