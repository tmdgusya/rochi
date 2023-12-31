package http

import (
	"errors"
	"fmt"
	"io"
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

func (res *Response) WriteTo(w io.Writer) (n int64, err error) {
	printf := func(format string, args ...any) error {
		m, err := fmt.Fprintf(w, format, args...)
		n += int64(m)
		return err
	}
	if err := printf("HTTP/1.1 %d %s\r\n", res.StatusCode, http.StatusText(res.StatusCode)); err != nil {
		return n, err
	}
	for _, h := range res.Headers {
		if err := printf("%s: %s\r\n", h.Key, h.Value); err != nil {
			return n, err
		}

	}
	if err := printf("\r\n%s\r\n", res.Body); err != nil {
		return n, err
	}
	return n, nil
}
