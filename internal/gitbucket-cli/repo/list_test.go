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

	// Test cases
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
			name:        "Successful list user repos",
			args:        []string{},
			expectedURL: "https://api.gitbucket.example.com/user/repos",
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
			name:        "Successful list org repos",
			args:        []string{"org-name"},
			expectedURL: "https://api.gitbucket.example.com/orgs/org-name/repos",
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
			args:          []string{},
			expectedURL:   "https://api.gitbucket.example.com/user/repos",
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to fetch repos"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name:        "Presenter error",
			args:        []string{},
			expectedURL: "https://api.gitbucket.example.com/user/repos",
			mockHTTPResp: &http.Response{
				StatusCode: http.StatusOK,
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

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Get(tt.expectedURL).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the List function
			err := repo.List(cmd, tt.args, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
