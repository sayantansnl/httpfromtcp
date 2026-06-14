package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])
	val := string(parts[1])

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("invalid header name: %q", key)
	}

	h.Set(key, strings.TrimSpace(val))
	return idx + 2, false, nil
}

func (h Headers) Set(key, val string) {
	h[key] = val
}
