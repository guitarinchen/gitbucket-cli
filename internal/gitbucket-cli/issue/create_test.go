package issue_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/issue"
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
	configs.UserName = "user"

	// Define test cases
	tests := []struct {
		name            string
		flags           issue.PossibleCreateFlags
		expectedURL     string
		expectedPayload issue.CreatePayload
		mockHTTPResp    *http.Response
		mockHTTPError   error
		mockPresenter   func()
		expectedError   bool
	}{
		{
			name: "Successful create issue with org/repo",
			flags: issue.PossibleCreateFlags{
				Repo:      "org/repo",
				Title:     "Issue Title",
				Body:      "Issue Body",
				Assignees: []string{"assignee1"},
				Milestone: 0,
				Labels:    []string{"bug"},
			},
			expectedURL: "https://api.gitbucket.example.com/repos/org/repo/issues",
			expectedPayload: issue.CreatePayload{
				Title:     "Issue Title",
				Body:      "Issue Body",
				Assignees: []string{"assignee1"},
				Milestone: 0,
				Labels:    []string{"bug"},
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
			name: "Successful create issue with user repo",
			flags: issue.PossibleCreateFlags{
				Repo:      "repo",
				Title:     "Another Issue",
				Body:      "Detailed Body",
				Assignees: nil,
				Milestone: 0,
				Labels:    nil,
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/issues",
			expectedPayload: issue.CreatePayload{
				Title:     "Another Issue",
				Body:      "Detailed Body",
				Assignees: nil,
				Milestone: 0,
				Labels:    nil,
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
			name: "HTTP client error",
			flags: issue.PossibleCreateFlags{
				Repo:  "repo",
				Title: "Issue with error",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/issues",
			expectedPayload: issue.CreatePayload{
				Title: "Issue with error",
			},
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to create issue"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name: "Presenter error",
			flags: issue.PossibleCreateFlags{
				Repo:  "repo",
				Title: "Issue with presenter error",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/issues",
			expectedPayload: issue.CreatePayload{
				Title: "Issue with presenter error",
			},
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(errors.New("presenter error"))
			},
			expectedError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Set up flags
			issue.CreateFlags = &tt.flags

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Post(tt.expectedURL, tt.expectedPayload).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the Create function
			err := issue.Create(cmd, nil, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
