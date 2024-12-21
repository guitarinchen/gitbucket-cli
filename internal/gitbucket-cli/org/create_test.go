package org_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/configs"
	"github.com/guitarinchen/gitbucket-cli/internal/gitbucket-cli/org"
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

	// Define test cases
	tests := []struct {
		name          string
		createFlags   *org.PossibleCreateFlags
		expectedURL   string
		mockHTTPResp  *http.Response
		mockHTTPError error
		mockPresenter func()
		expectedError bool
	}{
		{
			name: "Successful create organization",
			createFlags: &org.PossibleCreateFlags{
				Login:       "new-org",
				Admin:       "admin-user",
				ProfileName: "New Organization",
				URL:         "https://example.com",
			},
			expectedURL: "https://api.gitbucket.example.com/admin/organizations",
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
			name: "HTTP client error",
			createFlags: &org.PossibleCreateFlags{
				Login:       "new-org",
				Admin:       "admin-user",
				ProfileName: "New Organization",
				URL:         "https://example.com",
			},
			expectedURL:   "https://api.gitbucket.example.com/admin/organizations",
			mockHTTPResp:  nil,
			mockHTTPError: errors.New("failed to create organization"),
			mockPresenter: nil,
			expectedError: true,
		},
		{
			name: "Presenter error",
			createFlags: &org.PossibleCreateFlags{
				Login:       "new-org",
				Admin:       "admin-user",
				ProfileName: "New Organization",
				URL:         "https://example.com",
			},
			expectedURL: "https://api.gitbucket.example.com/admin/organizations",
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

			// Set CreateFlags for the test
			org.CreateFlags = tt.createFlags

			// Set up HTTP client behavior
			if tt.mockHTTPResp != nil || tt.mockHTTPError != nil {
				mockHTTPClient.EXPECT().Post(tt.expectedURL, gomock.Any()).Return(tt.mockHTTPResp, tt.mockHTTPError)
			}

			// Set up Presenter behavior
			if tt.mockPresenter != nil {
				tt.mockPresenter()
			}

			// Run the Create function
			err := org.Create(cmd, []string{}, mockHTTPClient, presenterFactory)

			// Verify the result
			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
