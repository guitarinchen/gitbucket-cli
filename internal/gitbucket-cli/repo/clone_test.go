package repo_test

import (
	"errors"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/repo"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
)

// MockCmd は、exec.Command を模倣するモック
type MockCmd struct {
	mock.Mock
}

func (m *MockCmd) Run() error {
	args := m.Called()
	return args.Error(0)
}

func TestClone(t *testing.T) {
	// Save the original Command
	originalCommand := repo.Command
	defer func() { repo.Command = originalCommand }()

	// Mock configurations
	configs.ApiBaseURL = "https://gitbucket.example.com"
	configs.UserName = "testuser"

	tests := []struct {
		name          string
		args          []string
		flags         *repo.PossibleCloneFlags
		expectedURL   string
		mockRunError  error
		expectedError bool
	}{
		{
			name: "Clone user repository with HTTPS",
			args: []string{"repo-name"},
			flags: &repo.PossibleCloneFlags{
				UseSSH: false,
				Branch: "main",
			},
			expectedURL:   "https://gitbucket.example.com/git/testuser/repo-name.git",
			mockRunError:  nil,
			expectedError: false,
		},
		{
			name: "Clone organization repository with SSH",
			args: []string{"org-name/repo-name"},
			flags: &repo.PossibleCloneFlags{
				UseSSH: true,
				Branch: "develop",
			},
			expectedURL:   "ssh://git@gitbucket.example.com:29418/org-name/repo-name.git",
			mockRunError:  nil,
			expectedError: false,
		},
		{
			name: "Error in git command",
			args: []string{"repo-name"},
			flags: &repo.PossibleCloneFlags{
				UseSSH: false,
				Branch: "main",
			},
			expectedURL:   "https://gitbucket.example.com/git/testuser/repo-name.git",
			mockRunError:  errors.New("git command failed"),
			expectedError: true,
		},
		{
			name: "Empty repository name",
			args: []string{""},
			flags: &repo.PossibleCloneFlags{
				UseSSH: false,
				Branch: "main",
			},
			expectedURL:   "",
			mockRunError:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the CloneFlags
			repo.CloneFlags = tt.flags

			// Mock exec.Command
			mockCmd := new(MockCmd)
			mockCmd.On("Run").Return(tt.mockRunError)
			repo.Command = func(name string, arg ...string) repo.Runner {
				if len(arg) > 0 && arg[0] == "-b" {
					if arg[1] != tt.flags.Branch || arg[2] != tt.expectedURL {
						t.Errorf("expected branch: %s and URL: %s, got branch: %s and URL: %s", tt.flags.Branch, tt.expectedURL, arg[1], arg[2])
					}
				}
				return mockCmd
			}

			// Execute the Clone function
			cmd := &cobra.Command{}
			err := repo.Clone(cmd, tt.args)

			// Verify expectations
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			if err == nil {
				mockCmd.AssertCalled(t, "Run")
			}
		})
	}
}
