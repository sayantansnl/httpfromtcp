package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/sayantansnl/httpfromtcp/internal/request"
	"github.com/sayantansnl/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

// Server is an HTTP 1.1 server
type Server struct {
	handler  Handler
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		handler:  handler,
		listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	rw := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		rw.WriteStatusLine(response.StatusCodeBadRequest)
		body := []byte(fmt.Sprintf("error in parsing request: %v", err))
		headers := response.GetDefaultHeaders(len(body))
		rw.WriteHeaders(headers)
		rw.WriteBody(body)
	}

	s.handler(rw, req)
}
