package http

import (
	"strings"
)

// Header represents an HTTP header. An HTTP header is a key-value pair, separated by a colon(:);
// The ky should be formatted in Title-Case.
// Use Request.AddHeader() or Response.AddHeader() to add headers to a request or response
// and guarantee title-casing of the key.
type Header struct {
	Key, Value string
}

// AsTitle returns the given header key as title case; e.g. "content-type" -> "Content-Type"
// You can implement this to use the standard library in Go
// see https://pkg.go.dev/net/textproto#CanonicalMIMEHeaderKey
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
