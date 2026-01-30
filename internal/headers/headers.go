package headers

import (
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
	if len(data) >= 2 && data[0] == '\r' && data[1] == '\n' {
		return 2, true, nil
	}

	crlfIndex := -1

	for i := 0; i+1 < len(data); i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			crlfIndex = i
			break
		}
	}

	if crlfIndex == -1 {
		return 0, false, nil
	}

	fieldLine := strings.Trim(string(data[:crlfIndex]), " ")

	for i := range fieldLine {
		if fieldLine[i] == ':' && fieldLine[i-1] == ' ' {
			return 0, false, fmt.Errorf("error: invalid spacing")
		}
	}

	fieldLineParts := strings.Split(fieldLine, " ")
	fieldName := strings.Split(fieldLineParts[0], ":")[0]
	fieldValue := strings.Trim(strings.Split(fieldLineParts[1], crlf)[0], " ")

	h[fieldName] = fieldValue

	return crlfIndex + 2, false, nil
}
