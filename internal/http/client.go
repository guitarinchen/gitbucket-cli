//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE
package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	net_http "net/http"
)

type HTTPClient struct {
	net_http.Client
	token string
}

type Client interface {
	Get(url string) (*net_http.Response, error)
	Post(url string, payload interface{}) (*net_http.Response, error)
}

func NewClient(token string) *HTTPClient {
	return &HTTPClient{
		token: token,
	}
}

func convertPayloadToBody(payload interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %v", err)
	}
	return bytes.NewBuffer(jsonData), nil
}

func (c *HTTPClient) sendRequest(method, url string, body io.Reader) (*net_http.Response, error) {
	req, err := net_http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %v", err)
	}

	// Add headers
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	}
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return resp, fmt.Errorf("send request: %v", err)
	}

	return resp, nil
}

func (c *HTTPClient) Get(url string) (*net_http.Response, error) {
	return c.sendRequest("GET", url, nil)
}

func (c *HTTPClient) Post(url string, payload interface{}) (*net_http.Response, error) {
	body, err := convertPayloadToBody(payload)
	if err != nil {
		return nil, err
	}
	return c.sendRequest("POST", url, body)
}
