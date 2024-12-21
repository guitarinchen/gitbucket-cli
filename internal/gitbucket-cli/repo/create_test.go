package repo_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/repo"
	"github.com/guitarinchen/gitbucket-cli/internal/http/mock"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter/mock"
	"github.com/spf13/cobra"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHTTPClient := mock_http.NewMockClient(ctrl)
	mockPresenter := mock_presenter.NewMockPresenter(ctrl)

	// Define presenter factory
	presenterFactory := func(resp *http.Response) presenter.Presenter {
		return mockPresenter
	}

	configs.ApiBaseURL = "https://api.gitbucket.example.com"

	// Test cases
	tests := []struct {
		name            string
		args            []string
		flags           *repo.PossibleCreateFlags
		expectedURL     string
		expectedPayload interface{}
		mockHTTPResp    *http.Response
		mockHTTPError   error
		mockPresenter   func()
		expectedError   bool
	}{
		{
			name:        "Successful create user repo",
			args:        []string{"my-repo"},
			flags:       &repo.PossibleCreateFlags{Description: "A user repo", IsPrivate: false, AutoInit: true},
			expectedURL: "https://api.gitbucket.example.com/user/repos",
			expectedPayload: repo.CreatePayload{
				Name:        "my-repo",
				Description: "A user repo",
				Private:     false,
				AutoInit:    true,
			},
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(nil)
			},
			expectedError: false,
		},
		{
			name:        "Successful create org repo",
			args:        []string{"org/my-repo"},
			flags:       &repo.PossibleCreateFlags{Description: "An org repo", IsPrivate: true, AutoInit: false},
			expectedURL: "https://api.gitbucket.example.com/orgs/org/repos",
			expectedPayload: repo.CreatePayload{
				Name:        "my-repo",
				Description: "An org repo",
				Private:     true,
				AutoInit:    false,
			},
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(nil)
			},
			expectedError: false,
		},
		{
			name:        "HTTP client error",
			args:        []string{"my-repo"},
			flags:       &repo.PossibleCreateFlags{Description: "A user repo", IsPrivate: false, AutoInit: true},
			expectedURL: "https://api.gitbucket.example.com/user/repos",
			expectedPayload: repo.CreatePayload{
				Name:        "my-repo",
				Description: "A user repo",
				Private:     false,
				AutoInit:    true,
			},
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to create repo"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name:        "Presenter error",
			args:        []string{"my-repo"},
			flags:       &repo.PossibleCreateFlags{Description: "A user repo", IsPrivate: false, AutoInit: true},
			expectedURL: "https://api.gitbucket.example.com/user/repos",
			expectedPayload: repo.CreatePayload{
				Name:        "my-repo",
				Description: "A user repo",
				Private:     false,
				AutoInit:    true,
			},
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(errors.New("presenter failed"))
			},
			expectedError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			repo.CreateFlags = tt.flags

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Post(tt.expectedURL, tt.expectedPayload).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the Create function
			err := repo.Create(cmd, tt.args, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
