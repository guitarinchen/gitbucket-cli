package issue

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
	Title     string   `json:"title"`
	Body      string   `json:"body,omitempty"`
	Assignees []string `json:"assignees"`
	Milestone int      `json:"milestone,omitempty"`
	Labels    []string `json:"labels"`
}

type PossibleCreateFlags struct {
	Repo      string
	Title     string
	Body      string
	Assignees []string
	Milestone int
	Labels    []string
}

var CreateFlags = &PossibleCreateFlags{}

func Create(cmd *cobra.Command, args []string, httpClient http.Client, presenterFactory func(resp *net_http.Response) presenter.Presenter) error {
	var url string
	if strings.Contains(CreateFlags.Repo, "/") {
		parts := strings.SplitN(CreateFlags.Repo, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/repos/%s/%s/issues", configs.ApiBaseURL, org, repo)
	} else {
		url = fmt.Sprintf("%s/repos/%s/%s/issues", configs.ApiBaseURL, configs.UserName, CreateFlags.Repo)
	}

	payload := CreatePayload{
		Title:     CreateFlags.Title,
		Body:      CreateFlags.Body,
		Assignees: CreateFlags.Assignees,
		Milestone: CreateFlags.Milestone,
		Labels:    CreateFlags.Labels,
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
