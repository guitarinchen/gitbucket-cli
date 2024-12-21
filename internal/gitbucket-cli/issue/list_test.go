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

func TestList(t *testing.T) {
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
		name          string
		args          []string
		expectedURL   string
		mockHTTPResp  *http.Response
		mockHTTPError error
		mockPresenter func()
		expectedError bool
	}{
		{
			name:        "Successful list issues with org/repo",
			args:        []string{"org/repo"},
			expectedURL: "https://api.gitbucket.example.com/repos/org/repo/issues?state=all",
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(nil)
			},
			expectedError: false,
		},
		{
			name:        "Successful list issues with user repo",
			args:        []string{"repo"},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/issues?state=all",
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			},
			mockHTTPError: nil,
			mockPresenter: func() {
				mockPresenter.EXPECT().Color().Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "HTTP client error",
			args:          []string{"org/repo"},
			expectedURL:   "https://api.gitbucket.example.com/repos/org/repo/issues?state=all",
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to fetch issues"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name:          "Empty repository argument",
			args:          []string{},
			expectedURL:   "",
			mockHTTPResp:  nil,
			mockHTTPError: nil,
			mockPresenter: nil,
			expectedError: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Get(tt.expectedURL).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the List function
			err := issue.List(cmd, tt.args, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
