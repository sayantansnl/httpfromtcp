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

type Writer struct {
	ResWriter io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	switch statusCode {
	case StatusOK:
		_, err := w.ResWriter.Write([]byte("HTTP/1.1 200 OK\r\n"))
		if err != nil {
			return fmt.Errorf("error in writing ok status line, error: %w", err)
		}

	case StatusBadRequest:
		_, err := w.ResWriter.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		if err != nil {
			return fmt.Errorf("error in writing bad request status line: %w", err)
		}

	case StatusServerError:
		_, err := w.ResWriter.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
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
	headers.Set("Content-Type", "text/html")

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for key, val := range headers {
		formattedHeader := fmt.Sprintf("%s: %s\r\n", key, val)
		_, err := w.ResWriter.Write([]byte(formattedHeader))
		if err != nil {
			return fmt.Errorf("unable to write header %s, error: %w", key, err)
		}
	}

	_, err := w.ResWriter.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.ResWriter.Write(p)
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkLength := len(p)
	nTotal := 0

	n, err := fmt.Fprintf(w.ResWriter, "%x\r\n", chunkLength)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.ResWriter.Write(p)
	if err != nil {
		return nTotal, err
	}

	nTotal += n

	n, err = w.ResWriter.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}

	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.ResWriter.Write([]byte("0\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for key, val := range h {
		_, err := w.ResWriter.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}

	_, err := w.ResWriter.Write([]byte("\r\n"))
	return err
}
