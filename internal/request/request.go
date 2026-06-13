package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	State       State
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type State int

const (
	initialized State = iota
	done
)

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	bytesRead, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqLine, err := parseRequestLine(bytesRead)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(bytesRead []byte) (*RequestLine, error) {
	crlfIdx := bytes.Index(bytesRead, []byte(crlf))
	if crlfIdx == -1 {
		return nil, fmt.Errorf("couldn't find CRLF in request line")
	}

	requestPart := strings.Split(string(bytesRead), crlf)[0]
	requestLine, err := requestLineFromRequestPart(requestPart)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromRequestPart(requestPart string) (*RequestLine, error) {
	requestLineParts := strings.Split(requestPart, " ")

	if len(requestLineParts) != 3 {
		return nil, fmt.Errorf("not enough request parts")
	}

	method := requestLineParts[0]
	for i := 0; i < len(method); i++ {
		if method[i] < 'A' || method[i] > 'Z' {
			return nil, fmt.Errorf("method should only be in capital letters")
		}
	}

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

	requestLine := RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: reqTarget,
		Method:        method,
	}

	return &requestLine, nil
}
