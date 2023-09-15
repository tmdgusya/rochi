package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Header represents an HTTP header. An HTTP header is a key-value pair, separated by a colon(:);
// The ky should be formatted in Title-Case.
// Use Request.AddHeader() or Response.AddHeader() to add headers to a request or response
// and guarantee title-casing of the key.
type Header struct {
	Key, Value string
}

// Request represents a HTTP 1.1 request.
type Request struct {
	Method  string
	Path    string
	Headers []Header
	Body    string // e.b, <html><body><h1>Hello, World!</h1></body></html>
}

// Response represents a HTTP Response
type Response struct {
	StatusCode int // e.g 200
	Headers    []Header
	Body       string
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

func (r *Request) WithHeader(key, value string) *Request {
	r.Headers = append(r.Headers, Header{AsTitle(key), value})
	return r
}

// AsTitle returns the given header key as title case; e.g. "content-type" -> "Content-Type"
func AsTitle(key string) string {
	if key == "" {
		panic("empty header key")
	}

	if isTitleCase(key) {
		return key
	}

	return newTitleCase(key)
}

func newTitleCase(key string) string {
	var b strings.Builder
	b.Grow(len(key))
	for i := range key {

		if i == 0 || key[i-1] == '-' {
			b.WriteByte(upper(key[i]))
			continue
		}
		b.WriteByte(lower(key[i]))
	}
	return b.String()
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 'a' - 'A'
	}
	return c
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c + 'A' - 'a'
	}
	return c
}

// isTitleCase returns true if the given header ky is already title case
func isTitleCase(key string) bool {

	for i := range key {
		if i == 0 || key[i-1] == '-' {
			// return false if the first character of the key is not upper-case
			if key[i] >= 'a' && key[i] <= 'z' {
				return false
			}
		} else if key[i] >= 'A' && key[i] <= 'Z' {
			// return false if the remain characters except for the first one of the key is not lower-case
			return false
		}
	}
	return true
}
