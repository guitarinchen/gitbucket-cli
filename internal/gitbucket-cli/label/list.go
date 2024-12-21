package label

import (
	"errors"
	"fmt"
	net_http "net/http"
	"strings"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/spf13/cobra"
)

func List(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	repo := args[0]
	if repo == "" {
		return errors.New("repository name is required")
	}

	var url string
	if strings.Contains(repo, "/") {
		parts := strings.SplitN(repo, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/repos/%s/%s/labels", configs.ApiBaseURL, org, repo)
	} else {
		url = fmt.Sprintf("%s/repos/%s/%s/labels", configs.ApiBaseURL, configs.UserName, repo)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	presenter := presenterFactory(resp)
	err = presenter.Color()
	if err != nil {
		return fmt.Errorf("error presenting response: %v", err)
	}

	return nil
}
