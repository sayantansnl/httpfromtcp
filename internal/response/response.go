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

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusCodeSuccess:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return fmt.Errorf("StatusCode: %d, error in writing: %w", StatusCodeSuccess, err)
		}
		return nil
	case StatusCodeBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return fmt.Errorf("StatusCode: %d, error in writing: %w", StatusCodeBadRequest, err)
		}
		return nil
	case StatusCodeServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return fmt.Errorf("StatusCode: %d, error in writing: %w", StatusCodeServerError, err)
		}
		return nil
	default:
		_, err := w.Write([]byte("HTTP/1.1 <code>\r\n"))
		return fmt.Errorf("StatusCode: none, error in writing: %w", err)
	}
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
