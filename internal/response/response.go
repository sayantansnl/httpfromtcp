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

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer:      w,
	}
}

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

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write in writer state %d", w.writerState)
	}

	defer func() {
		w.writerState = writerStateHeaders
	}()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()

	headers.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	headers.Set("Connection", "close")
	headers.Set("Content-Type", "text/plain")

	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write in writer state %d", w.writerState)
	}

	defer func() {
		w.writerState = writerStateBody
	}()

	for key, val := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write in writer state %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write in writer state: %d", w.writerState)
	}
	chunkSize := len(p)
	total := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	total += n
	if err != nil {
		return total, fmt.Errorf("unable to format, error: %w", err)
	}

	x, err := w.writer.Write(p)
	total += x
	if err != nil {
		return total, fmt.Errorf("unable to write in bytes, error: %w", err)
	}

	z, err := w.writer.Write([]byte("\r\n"))
	total += z
	if err != nil {
		return total, fmt.Errorf("unable to write crlf, error: %w", err)
	}
	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("unable to write body in writer state %d", w.writerState)
	}

	n, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return 0, fmt.Errorf("unable to write terminating sequence, error: %w", err)
	}

	return n, nil
}

func (w *Writer) WriteTrailers(headers headers.Headers) error {
	for key, value := range headers {
		header := fmt.Sprintf("%s: %s\r\n", key, value)
		if _, err := w.writer.Write([]byte(header)); err != nil {
			return err
		}
	}

	if _, err := w.writer.Write([]byte("\r\n")); err != nil {
		return err
	}

	return nil
}
