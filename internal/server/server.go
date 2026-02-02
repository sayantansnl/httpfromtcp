package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/sayantansnl/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portString)
	if err != nil {
		return nil, fmt.Errorf("unable to listen on port %s, error: %w", portString, err)
	}
	server := &Server{
		listener: l,
	}
	server.closed.Store(false)

	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.closed.Store(true)
			break
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)

	if err := response.WriteStatusLine(conn, 200); err != nil {
		log.Fatalf("cannot write status line to connection, error: %v", err)
	}

	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Fatalf("cannot write headers to connection, error: %v", err)
	}
}
