package presenter_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/fatih/color"
	"github.com/guitarinchen/gitbucket-cli/internal/presenter"
	"github.com/stretchr/testify/assert"
)

func TestResponsePresenter_Color(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedOutput string
		wantErr        bool
	}{
		{
			name:       "Successful response",
			statusCode: http.StatusOK,
			responseBody: `{
				"key": "value"
			}`,
			expectedOutput: `{
  "data": {
    "key": "value"
  }
}`,
			wantErr: false,
		},
		{
			name:       "Error response with valid JSON",
			statusCode: http.StatusBadRequest,
			responseBody: `{
				"message": "invalid request"
			}`,
			expectedOutput: `{
  "error": {
    "code": 400,
    "message": "invalid request"
  }
}`,
			wantErr: false,
		},
		{
			name:           "Error response with invalid JSON",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `invalid-json`,
			expectedOutput: "",
			wantErr:        true,
		},
		{
			name:         "Empty response body",
			statusCode:   http.StatusNoContent,
			responseBody: ``,
			expectedOutput: `{
  "data": null
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock HTTP response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(tt.responseBody)),
			}

			// Create a new ResponsePresenter
			presenter := presenter.NewResponsePresenter(resp)

			// Capture the colored output
			outputBuffer := &bytes.Buffer{}
			color.Output = outputBuffer // Override default output to capture printed output

			// Run the method
			err := presenter.Color()

			// Verify error
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, outputBuffer.String(), tt.expectedOutput)
			}
		})
	}
}
