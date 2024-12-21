package org

import (
	"fmt"
	net_http "net/http"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

type CreatePayload struct {
	Login       string `json:"login"`
	Admin       string `json:"admin"`
	ProfileName string `json:"profile_name"`
	Url         string `json:"url"`
}

type PossibleCreateFlags struct {
	Login       string
	Admin       string
	ProfileName string
	URL         string
}

var CreateFlags = &PossibleCreateFlags{}

func Create(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	url := fmt.Sprintf("%s/admin/organizations", configs.ApiBaseURL)

	payload := CreatePayload{
		Login:       CreateFlags.Login,
		Admin:       CreateFlags.Admin,
		ProfileName: CreateFlags.ProfileName,
		Url:         CreateFlags.URL,
	}

	resp, err := httpClient.Post(url, payload)
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
