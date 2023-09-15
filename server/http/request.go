package http

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Request represents a HTTP 1.1 request.
type Request struct {
	Method  string
	Path    string
	Headers []Header
	Body    string // e.b, <html><body><h1>Hello, World!</h1></body></html>
}

// NewRequest Create New Request instance with the following arguments
func NewRequest(method, path, host, body string) (*Request, error) {
	switch {
	case method == "":
		return nil, errors.New("missing required argument: method")
	case path == "":
		return nil, errors.New("missing required argument: path")
	case !strings.HasPrefix(path, "/"):
		return nil, errors.New("path must start with '/'")
	case host == "":
		return nil, errors.New("missing required argument: host")
	default:
		headers := make([]Header, 2)
		headers[0] = Header{Key: "Host", Value: host}
		if body != "" {
			headers = append(headers, Header{"Content-Length", fmt.Sprintf("%d", len(body))})
		}
		return &Request{
			Method:  method,
			Path:    path,
			Headers: headers,
			Body:    body,
		}, nil
	}
}

func (r *Request) WithHeader(key, value string) *Request {
	r.Headers = append(r.Headers, Header{AsTitle(key), value})
	return r
}

func (r *Request) WriteTo(w io.Writer) (n int64, err error) {
	// write & count bytes written
	// using small closures like this to cut down on repetition
	printf := func(format string, args ...any) error {
		// m is number of bytes written
		m, err := fmt.Fprintf(w, format, args...)
		n += int64(m)
		return err
	}

	if err := printf("%s %s HTTP/1.1\r\n", r.Method, r.Path); err != nil {
		return n, err
	}

	for _, h := range r.Headers {
		if err := printf("%s %s HTTP/1.1\r\n", h.Key, h.Value); err != nil {
			return n, err
		}
	}

	printf("\r\n")                 // Write the empty line that separates the headers from the body
	err = printf("%s\r\n", r.Body) // Write the body and terminate with a newline
	return n, err
}
