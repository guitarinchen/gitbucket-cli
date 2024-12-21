/*
Copyright © 2024 guitarinchen <guitarinchen@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitbucket-cli",
	Short: "Work seamlessly with GitBucket from the command line.",
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

func getGitConfigValue(section, key string) (string, error) {
	paths := []string{
		expandPath("~/.gitconfig"),
		expandPath("~/.config/git/config"),
	}

	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		inSection := false

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				// New section
				inSection = line == "["+section+"]"
				continue
			}

			if inSection {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[0]) == key {
					s := strings.TrimSpace(parts[1])
					if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
						return s[1 : len(s)-1], err
					}
					return s, nil
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("key not found")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	configs.ApiToken = os.Getenv("GITBUCKET_API_TOKEN")
	if configs.ApiToken == "" {
		println("GITBUCKET_API_TOKEN is not set.")
		os.Exit(1)
	}

	configs.ApiBaseURL = os.Getenv("GITBUCKET_API_BASE_URL")
	if configs.ApiBaseURL == "" {
		println("GITBUCKET_API_BASE_URL is not set.")
		os.Exit(1)
	}

	userName, err := getGitConfigValue("user", "name")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	configs.UserName = userName

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
