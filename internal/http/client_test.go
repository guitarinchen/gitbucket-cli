package http_test

import (
	"io"
	net_http "net/http"
	"net/http/httptest"
	"testing"

	"github.com/guitarinchen/gitbucket-cli/internal/http"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_Get(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		mockServer func(w net_http.ResponseWriter, r *net_http.Request)
		url        string
		wantStatus int
		wantErr    bool
	}{
		{
			name:  "Successful GET request",
			token: "valid_token",
			mockServer: func(w net_http.ResponseWriter, r *net_http.Request) {
				assert.Equal(t, "token valid_token", r.Header.Get("Authorization"))
				w.WriteHeader(net_http.StatusOK)
			},
			url:        "/test",
			wantStatus: net_http.StatusOK,
			wantErr:    false,
		},
		{
			name:  "GET request with 500 status code",
			token: "valid_token",
			mockServer: func(w net_http.ResponseWriter, r *net_http.Request) {
				w.WriteHeader(net_http.StatusInternalServerError)
			},
			url:        "/error",
			wantStatus: net_http.StatusInternalServerError,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server
			server := httptest.NewServer(net_http.HandlerFunc(tt.mockServer))
			defer server.Close()

			client := http.NewClient(tt.token)
			resp, err := client.Get(server.URL + tt.url)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestHTTPClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		token      string
		mockServer func(w net_http.ResponseWriter, r *net_http.Request)
		url        string
		payload    interface{}
		wantStatus int
		wantErr    bool
	}{
		{
			name:  "Successful POST request",
			token: "valid_token",
			mockServer: func(w net_http.ResponseWriter, r *net_http.Request) {
				body, _ := io.ReadAll(r.Body)
				assert.JSONEq(t, `{"key":"value"}`, string(body))
				assert.Equal(t, "token valid_token", r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				w.WriteHeader(net_http.StatusCreated)
			},
			url:        "/test",
			payload:    map[string]string{"key": "value"},
			wantStatus: net_http.StatusCreated,
			wantErr:    false,
		},
		{
			name:  "POST request with invalid payload",
			token: "valid_token",
			mockServer: func(w net_http.ResponseWriter, r *net_http.Request) {
				w.WriteHeader(net_http.StatusBadRequest)
			},
			url:        "/invalid",
			payload:    make(chan int), // Invalid payload
			wantStatus: 0,              // No response expected
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server
			server := httptest.NewServer(net_http.HandlerFunc(tt.mockServer))
			defer server.Close()

			client := http.NewClient(tt.token)
			resp, err := client.Post(server.URL+tt.url, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
