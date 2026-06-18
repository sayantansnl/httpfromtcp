package server

import (
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/request"
	"github/sayantansnl/httpfromtcp/internal/response"
	"net"
	"strconv"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	Addr     string
	Listener net.Listener
	IsClosed *atomic.Bool
	Handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := ":" + strconv.Itoa(port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error in creating a listener: %w", err)
	}

	server := Server{
		Addr:     addr,
		Listener: listener,
		IsClosed: &atomic.Bool{},
		Handler:  handler,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.IsClosed.Store(true)

	if err := s.Listener.Close(); err != nil {
		return fmt.Errorf("unable to close listener, error: %w", err)
	}

	return nil
}

func (s *Server) listen() {
	if s.IsClosed.Load() {
		return
	}

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			s.IsClosed.Store(true)
			return
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	writer := response.Writer{
		ResWriter: conn,
	}
	req, err := request.RequestFromReader(conn)
	if err != nil {
		writer.WriteStatusLine(response.StatusServerError)
		body := []byte(fmt.Sprintf("error in parsing request, error: %v", err))
		headers := response.GetDefaultHeaders(len(body))
		writer.WriteHeaders(headers)
		writer.WriteBody(body)
	}

	s.Handler(&writer, req)
}
