package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sayantansnl/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       requestState
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := &Request{
		State:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}
	for req.State != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, fmt.Errorf("incomplete request: %w", err)
				}
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, fmt.Errorf("error in parsing data for the whole request line: %w", err)
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done == true {
			r.State = requestStateParsingBody
		}

		return n, nil
	case requestStateParsingBody:
		contentLengthHeader := r.Headers.Get("Content-Length")
		if contentLengthHeader == "" {
			// nothing to parse
			r.State = requestStateDone
			return 0, nil
		}
		contentLength, err := strconv.Atoi(contentLengthHeader)
		if err != nil {
			return 0, fmt.Errorf("invalid content length: %w", err)
		}

		for _, part := range data {
			r.Body = append(r.Body, part)
		}

		if len(r.Body) > contentLength {
			fmt.Println("body too much")
			return 0, fmt.Errorf("body is greater than the specified content length")
		}

		if len(r.Body) < contentLength {
			fmt.Println("body is shorter than content length")
			return len(data), nil
		}

		if len(r.Body) == contentLength {
			r.State = requestStateDone
		}

		return len(r.Body), nil
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
