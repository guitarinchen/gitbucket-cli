package org

import (
	"fmt"
	net_http "net/http"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

func List(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	url := fmt.Sprintf("%s/user/orgs", configs.ApiBaseURL)

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("send request: %v", err)
	}

	presenter := presenterFactory(resp)
	err = presenter.Color()
	if err != nil {
		return fmt.Errorf("present response: %v", err)
	}

	return nil
}
