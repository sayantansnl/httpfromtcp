package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestPart := strings.Split(string(bytes), "\r\n")[0]

	reqLine, err := parseRequestLine(requestPart)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(requestPart string) (*RequestLine, error) {
	requestLineParts := strings.Split(requestPart, " ")

	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("not enough request parts")
	}

	method := requestLineParts[0]
	reqTarget := requestLineParts[1]

	httpPart := strings.Split(requestLineParts[2], "/")

	if len(httpPart) != 2 {
		return nil, fmt.Errorf("unrecognized http")
	}
	http := httpPart[0]
	httpVersion := httpPart[1]

	if http != "HTTP" {
		return nil, fmt.Errorf("unrecognized protocol")
	}

	if httpVersion != "1.1" {
		return nil, fmt.Errorf("only version 1.1 is supported")
	}

	for i := 0; i < len(method); i++ {
		if method[i] < 'A' || method[i] > 'Z' {
			return nil, fmt.Errorf("method should only be in capital letters")
		}
	}

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        method,
	}

	return &requestLine, nil
}
