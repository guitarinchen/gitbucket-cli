package repo

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

type CreatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

type PossibleCreateFlags struct {
	IsPrivate   bool
	Description string
	AutoInit    bool
}

var CreateFlags = &PossibleCreateFlags{}

func Create(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	repoName := args[0]
	if repoName == "" {
		return errors.New("empty repository")
	}

	var url string
	var payload CreatePayload

	if strings.Contains(repoName, "/") {
		parts := strings.SplitN(repoName, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/orgs/%s/repos", configs.ApiBaseURL, org)
		payload = CreatePayload{
			Name:        repo,
			Description: CreateFlags.Description,
			Private:     CreateFlags.IsPrivate,
			AutoInit:    CreateFlags.AutoInit,
		}
	} else {
		url = fmt.Sprintf("%s/user/repos", configs.ApiBaseURL)
		payload = CreatePayload{
			Name:        repoName,
			Description: CreateFlags.Description,
			Private:     CreateFlags.IsPrivate,
			AutoInit:    CreateFlags.AutoInit,
		}
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
