package issue

import (
	"errors"
	"fmt"
	net_http "net/http"
	net_url "net/url"
	"strings"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

type PossibleListFlags struct {
	State string
}

var ListFlags = &PossibleListFlags{
	State: "all",
}

func List(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	if len(args) < 1 || args[0] == "" {
		return errors.New("repository name is required")
	}

	repo := args[0]
	params := net_url.Values{}
	params.Add("state", ListFlags.State)

	var url string
	if strings.Contains(repo, "/") {
		parts := strings.SplitN(repo, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/repos/%s/%s/issues?%s", configs.ApiBaseURL, org, repo, params.Encode())
	} else {
		url = fmt.Sprintf("%s/repos/%s/%s/issues?%s", configs.ApiBaseURL, configs.UserName, repo, params.Encode())
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch issues from %s: %v", url, err)
	}

	presenter := presenterFactory(resp)
	if err := presenter.Color(); err != nil {
		return fmt.Errorf("present response: %v", err)
	}

	return nil
}
