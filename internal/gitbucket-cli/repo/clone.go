package repo

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	net_url "net/url"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/spf13/cobra"
)

type PossibleCloneFlags struct {
	UseSSH bool
	Branch string
}

var CloneFlags = &PossibleCloneFlags{
	Branch: "main",
}

type Runner interface {
	Run() error
}

var Command = func(name string, arg ...string) Runner {
	return exec.Command(name, arg...)
}

func Clone(cmd *cobra.Command, args []string) error {
	repoName := args[0]
	if repoName == "" {
		return errors.New("empty repository")
	}

	parsedUrl, err := net_url.Parse(configs.ApiBaseURL)
	if err != nil {
		return fmt.Errorf("parse url: %v", err)
	}
	domain := parsedUrl.Host
	var baseUrl string
	if CloneFlags.UseSSH {
		sshPort := 29418
		baseUrl = fmt.Sprintf("ssh://git@%s:%d", domain, sshPort)
	} else {
		baseUrl = fmt.Sprintf("%s/git", fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host))
	}

	var url string
	if strings.Contains(repoName, "/") {
		parts := strings.SplitN(repoName, "/", 2)
		org := parts[0]
		repo := parts[1]
		url = fmt.Sprintf("%s/%s/%s.git", baseUrl, org, repo)
	} else {
		url = fmt.Sprintf("%s/%s/%s.git", baseUrl, configs.UserName, repoName)
	}

	c := Command("git", "clone", "-b", CloneFlags.Branch, url)
	fmt.Printf("Cloning %s...\n", url)
	if err := c.Run(); err != nil {
		return fmt.Errorf("run command: %v", err)
	}

	fmt.Println("The repository has been successfully cloned.")

	return nil
}
