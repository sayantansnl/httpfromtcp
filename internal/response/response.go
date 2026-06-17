package response

import (
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOK          StatusCode = 200
	StatusBadRequest  StatusCode = 400
	StatusServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return fmt.Errorf("error in writing ok status line, error: %w", err)
		}

	case StatusBadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return fmt.Errorf("error in writing bad request status line: %w", err)
		}

	case StatusServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		if err != nil {
			return fmt.Errorf("error in writing server error status line, error: %w", err)
		}

	default:
		return fmt.Errorf("unknow status")
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", strconv.Itoa(contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		formattedHeader := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.Write([]byte(formattedHeader))
		if err != nil {
			return fmt.Errorf("unable to write header %s, error: %w", key, err)
		}
	}

	w.Write([]byte("\r\n"))
	return nil
}
