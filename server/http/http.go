package http

import (
	"bytes"
	"encoding"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Interface Section -- Start
var _, _ fmt.Stringer = (*Request)(nil), (*Response)(nil) // compile-time check that Request and Response implement fmt.Stringer
var _, _ encoding.TextMarshaler = (*Request)(nil), (*Response)(nil)

func (r *Request) String() string     { b := new(strings.Builder); r.WriteTo(b); return b.String() }
func (resp *Response) String() string { b := new(strings.Builder); resp.WriteTo(b); return b.String() }
func (r *Request) MarshalText() ([]byte, error) {
	b := new(bytes.Buffer)
	r.WriteTo(b)
	return b.Bytes(), nil
}
func (resp *Response) MarshalText() ([]byte, error) {
	b := new(bytes.Buffer)
	resp.WriteTo(b)
	return b.Bytes(), nil
}

// Interface Section -- End

// ParseRequest parses a HTTP request from the given text.
func ParseRequest(raw string) (r Request, err error) {
	// request has three parts:
	// 1. Request linedd
	// 2. Headers
	// 3. Body (optional)
	lines := splitLines(raw)

	log.Println(lines)
	if len(lines) < 3 {
		return Request{}, fmt.Errorf("malformed request: should have at least 3 lines")
	}
	// First line is special.
	first := strings.Fields(lines[0])
	r.Method, r.Path = first[0], first[1]
	if !strings.HasPrefix(r.Path, "/") {
		return Request{}, fmt.Errorf("malformed request: path should start with /")
	}
	if !strings.Contains(first[2], "HTTP") {
		return Request{}, fmt.Errorf("malformed request: first line should contain HTTP version")
	}
	var foundhost bool
	var bodyStart int
	// then we have headers, up until the an empty line.
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" { // empty line
			bodyStart = i + 1
			break
		}
		key, val, ok := strings.Cut(lines[i], ": ")
		if !ok {
			return Request{}, fmt.Errorf("malformed request: header %q should be of form 'key: value'", lines[i])
		}
		if key == "Host" { // special case: host header is required.
			foundhost = true
		}
		key = AsTitle(key)

		r.Headers = append(r.Headers, Header{key, val})
	}
	end := len(lines) - 1 // recombine the body using normal newlines; skip the last empty line.
	r.Body = strings.Join(lines[bodyStart:end], "\r\n")
	if !foundhost {
		return Request{}, fmt.Errorf("malformed request: missing Host header")
	}
	return r, nil
}

// ParseResponse parses the given HTTP/1.1 response string into the Response. It returns an error if the Response is invalid,
// - not a valid integer
// - invalid status code
// - missing status text
// - invalid headers
// it doesn't properly handle multi-line headers, headers with multiple values, or html-encoding, etc.zzs
func ParseResponse(raw string) (resp *Response, err error) {
	// response has three parts:
	// 1. Response line
	// 2. Headers
	// 3. Body (optional)
	lines := splitLines(raw)
	log.Println(lines)

	// First line is special.
	first := strings.SplitN(lines[0], " ", 3)
	if !strings.Contains(first[0], "HTTP") {
		return nil, fmt.Errorf("malformed response: first line should contain HTTP version")
	}
	resp = new(Response)
	resp.StatusCode, err = strconv.Atoi(first[1])
	if err != nil {
		return nil, fmt.Errorf("malformed response: expected status code to be an integer, got %q", first[1])
	}
	if first[2] == "" || http.StatusText(resp.StatusCode) != first[2] {
		log.Printf("missing or incorrect status text for status code %d: expected %q, but got %q", resp.StatusCode, http.StatusText(resp.StatusCode), first[2])
	}
	var bodyStart int
	// then we have headers, up until the an empty line.
	for i := 1; i < len(lines); i++ {
		log.Println(i, lines[i])
		if lines[i] == "" { // empty line
			bodyStart = i + 1
			break
		}
		key, val, ok := strings.Cut(lines[i], ": ")
		if !ok {
			return nil, fmt.Errorf("malformed response: header %q should be of form 'key: value'", lines[i])
		}
		key = AsTitle(key)
		resp.Headers = append(resp.Headers, Header{key, val})
	}
	resp.Body = strings.TrimSpace(strings.Join(lines[bodyStart:], "\r\n")) // recombine the body using normal newlines.
	return resp, nil
}

// splitLines on the "\r\n" sequence; multiple separators in a row are NOT collapsed.
func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var lines []string
	i := 0
	for {
		j := strings.Index(s[i:], "\r\n")
		if j == -1 {
			lines = append(lines, s[i:])
			return lines
		}
		lines = append(lines, s[i:i+j]) // up to but not including the \r\n
		i += j + 2                      // skip the \r\n
	}
}
