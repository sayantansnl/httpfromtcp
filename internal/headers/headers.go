package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func NewHeaders() Headers {
	headers := make(Headers)
	return headers
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, nil
	}
	if crlfIndex == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:crlfIndex], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(string(key), " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := string(bytes.TrimSpace(parts[1]))
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)

	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid header token has been in found: %s", key)
	}

	h.Set(key, value)
	return crlfIndex + 2, false, nil
}

func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value
}

func (h Headers) Get(key string) string {
	keyLowerCase := strings.ToLower(key)

	val, ok := h[keyLowerCase]
	if !ok {
		return ""
	}

	return val
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Delete(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}
