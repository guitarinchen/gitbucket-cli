package label

import (
	"fmt"
	net_http "net/http"
	"strings"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

type CreatePayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type PossibleCreateFlags struct {
	Repo  string
	Name  string
	Color string
}

var CreateFlags = &PossibleCreateFlags{}

func Create(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	var url string
	if strings.Contains(CreateFlags.Repo, "/") {
		parts := strings.SplitN(CreateFlags.Repo, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/repos/%s/%s/labels", configs.ApiBaseURL, org, repo)
	} else {
		url = fmt.Sprintf("%s/repos/%s/%s/labels", configs.ApiBaseURL, configs.UserName, CreateFlags.Repo)
	}

	payload := CreatePayload{
		Name:  CreateFlags.Name,
		Color: CreateFlags.Color,
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
