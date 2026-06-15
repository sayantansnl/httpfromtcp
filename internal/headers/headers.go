package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

var tokenChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return idx + 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])
	val := string(parts[1])

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("invalid header name: %q", key)
	}

	if err := ensureValidKey(key); err != nil {
		return 0, false, err
	}

	h.Set(key, strings.TrimSpace(val))
	return idx + 2, false, nil
}

func (h Headers) Set(key, val string) {
	lowerKey := strings.ToLower(key)
	lowerVal := strings.ToLower(val)

	if _, ok := h[lowerKey]; !ok {
		h[lowerKey] = lowerVal
		return
	}

	originalVal := h[lowerKey]
	splittedOriginalVal := strings.Split(originalVal, ", ")
	splittedOriginalVal = append(splittedOriginalVal, lowerVal)
	h[lowerKey] = strings.Join(splittedOriginalVal, ", ")
}

func ensureValidKey(key string) error {
	for _, k := range key {
		if !isValidToken(k) {
			return fmt.Errorf("invalid token in key")
		}
	}
	return nil
}

func isValidToken(b rune) bool {
	if b >= 'A' && b <= 'Z' ||
		b >= 'a' && b <= 'z' ||
		b >= '0' && b <= '9' {
		return true
	}

	return slices.Contains(tokenChars, b)
}
