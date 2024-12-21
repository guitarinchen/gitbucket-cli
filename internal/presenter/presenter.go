//go:generate mockgen -source=$GOFILE -destination=mock/mock_$GOFILE
package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
)

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Presenter interface {
	Color() error
}

type ResponsePresenter struct {
	Resp *http.Response
}

func NewResponsePresenter(resp *http.Response) *ResponsePresenter {
	return &ResponsePresenter{Resp: resp}
}

func (r *ResponsePresenter) Color() error {
	body, err := io.ReadAll(r.Resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %v", err)
	}
	defer r.Resp.Body.Close()

	// Handle empty body
	if len(body) == 0 {
		if r.Resp.StatusCode >= 200 && r.Resp.StatusCode < 300 {
			body = []byte(`null`)
		}
	}

	var result interface{}
	if r.Resp.StatusCode >= 200 && r.Resp.StatusCode < 300 {
		var data interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			return err
		}

		result = SuccessResponse{
			Data: data,
		}
	} else {
		var error Error
		if err := json.Unmarshal(body, &error); err != nil {
			return err
		}
		error.Code = r.Resp.StatusCode

		result = ErrorResponse{
			Error: error,
		}
	}

	formatted, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("format json: %v", err)
	}

	c := color.New(color.FgYellow)
	c.Println(string(formatted))

	return nil
}
