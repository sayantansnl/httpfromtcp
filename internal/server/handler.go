package server

import (
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/request"
	"github/sayantansnl/httpfromtcp/internal/response"
	"io"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func WriteError(w io.Writer, handlerError HandlerError) error {
	err := response.WriteStatusLine(w, handlerError.StatusCode)
	if err != nil {
		return fmt.Errorf("unable to write status line, error: %v", err)
	}

	headers := response.GetDefaultHeaders(len([]byte(handlerError.Message)))
	err = response.WriteHeaders(w, headers)
	if err != nil {
		return fmt.Errorf("unable to write headers, error: %w", err)
	}

	if _, err := w.Write([]byte(handlerError.Message)); err != nil {
		return fmt.Errorf("unable to write error message to io.Writer, error: %w", err)
	}

	return nil
}
