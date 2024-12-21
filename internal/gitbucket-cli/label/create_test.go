package label_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/label"
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
		flags           *label.PossibleCreateFlags
		expectedURL     string
		expectedPayload label.CreatePayload
		mockHTTPResp    *http.Response
		mockHTTPError   error
		mockPresenter   func()
		expectedError   bool
	}{
		{
			name: "Successful create label with org/repo",
			flags: &label.PossibleCreateFlags{
				Repo:  "org/repo",
				Name:  "bug",
				Color: "ff0000",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/org/repo/labels",
			expectedPayload: label.CreatePayload{
				Name:  "bug",
				Color: "ff0000",
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
			name: "Successful create label with user repo",
			flags: &label.PossibleCreateFlags{
				Repo:  "repo",
				Name:  "feature",
				Color: "00ff00",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/labels",
			expectedPayload: label.CreatePayload{
				Name:  "feature",
				Color: "00ff00",
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
			flags: &label.PossibleCreateFlags{
				Repo:  "repo",
				Name:  "invalid",
				Color: "123456",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/labels",
			expectedPayload: label.CreatePayload{
				Name:  "invalid",
				Color: "123456",
			},
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to create label"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name: "Presenter error",
			flags: &label.PossibleCreateFlags{
				Repo:  "repo",
				Name:  "presenter-error",
				Color: "abcdef",
			},
			expectedURL: "https://api.gitbucket.example.com/repos/user/repo/labels",
			expectedPayload: label.CreatePayload{
				Name:  "presenter-error",
				Color: "abcdef",
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

			// Set the flags for the test case
			label.CreateFlags = tt.flags

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Post(tt.expectedURL, tt.expectedPayload).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the Create function
			err := label.Create(cmd, nil, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
