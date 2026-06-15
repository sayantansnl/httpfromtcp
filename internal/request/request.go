package request

import (
	"bytes"
	"errors"
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       State
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type State int

const (
	requestStateInitialized State = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	req := Request{
		State:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.State != requestStateDone {
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		bytesRead, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, fmt.Errorf("incomplete request")
				}
			}
			return nil, err
		}

		readToIndex += bytesRead

		numBytes, err := req.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buff, buff[numBytes:])
		readToIndex -= numBytes
	}

	return &req, nil
}

func parseRequestLine(bytesRead []byte) (*RequestLine, int, error) {
	crlfIdx := bytes.Index(bytesRead, []byte(crlf))
	if crlfIdx == -1 {
		return nil, 0, nil
	}

	requestPart := string(bytesRead[:crlfIdx])
	requestLine, err := requestLineFromRequestPart(requestPart)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, crlfIdx + 2, nil
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

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 && err == nil {
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
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateParsingBody
			return n, nil
		}
		return n, err
	case requestStateParsingBody:
		contentLength, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.State = requestStateDone
			return len(data), nil
		}

		cl, err := strconv.Atoi(contentLength)
		if err != nil {
			return 0, fmt.Errorf("couldn't convert string to int, error: %w", err)
		}

		r.Body = append(r.Body, data...)
		bodyLength := len(r.Body)
		if bodyLength > cl {
			return 0, fmt.Errorf("provided content doesn't match content length")
		}
		if bodyLength == cl {
			r.State = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
