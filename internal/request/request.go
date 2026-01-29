package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
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
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot use io.ReadAll: %w", err)
	}

	req, err := parseRequestLine(string(b))
	if err != nil {
		return nil, fmt.Errorf("unable to parse: %w", err)
	}

	return &Request{
		RequestLine: *req,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	requestLine := strings.Split(line, "\r\n")[0]
	parts := strings.Split(requestLine, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of parts")
	}

	method := parts[0]
	target := parts[1]
	httpSlashVersion := parts[2]
	protocol := strings.Split(httpSlashVersion, "/")[0]
	version := strings.Split(httpSlashVersion, "/")[1]

	for _, r := range method {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return nil, fmt.Errorf("method should contain only capital alphabetic letters")
		}
	}

	if protocol != "HTTP" {
		return nil, fmt.Errorf("unsupported protocol")
	}

	if version != "1.1" {
		return nil, fmt.Errorf("only version 1.1 supported")
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}
