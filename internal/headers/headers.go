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

	header := string(data[:idx])

	key := strings.Split(header, ": ")[0]
	val := strings.Split(header, ": ")[1]

	valWithoutCRLF := strings.Split(val, crlf)[0]

	if key != strings.TrimSpace(key) {
		return 0, false, fmt.Errorf("poor formatting")
	}

	h.Set(key, strings.TrimSpace(valWithoutCRLF))

	h.Parse(data[idx+2:])

	n += idx + 2
	return n, false, nil
}

func (h Headers) Set(key, val string) {
	h[key] = val
}
