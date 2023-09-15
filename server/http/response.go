package http

import (
	"errors"
	"fmt"
	"net/http"
)

// Response represents a HTTP Response
type Response struct {
	StatusCode int // e.g 200
	Headers    []Header
	Body       string
}

// NewResponse create new Response instance with the following arguments
func NewResponse(status int, body string) (*Response, error) {
	switch {
	case status < 100 || status > 599:
		return nil, errors.New("invalid status code")
	default:
		if body == "" {
			body = http.StatusText(status)
		}
		headers := make([]Header, 1)
		headers[0] = Header{
			"Content-Length",
			fmt.Sprintf("%d", len(body)),
		}
		return &Response{
			StatusCode: status,
			Headers:    headers,
			Body:       body,
		}, nil
	}
}

func (res *Response) WithHeader(key, value string) *Response {
	res.Headers = append(res.Headers, Header{AsTitle(key), value})
	return res
}
