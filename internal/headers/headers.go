package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	headers := make(Headers)
	return headers
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := bytes.Index(data, []byte(crlf))
	if crlfIndex == -1 {
		return 0, false, fmt.Errorf("not enough data provided")
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

	h[key] = value
	return crlfIndex + 2, false, nil
}
