package response

import (
	"fmt"
	"io"

	"github.com/sayantansnl/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeSuccess     StatusCode = 200
	StatusCodeBadRequest  StatusCode = 400
	StatusCodeServerError StatusCode = 500
)

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
