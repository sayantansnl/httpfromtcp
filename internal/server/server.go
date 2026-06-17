package server

import (
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	Addr     string
	Listener net.Listener
	IsClosed *atomic.Bool
}

func Serve(port int) (*Server, error) {
	addr := ":" + strconv.Itoa(port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error in creating a listener: %w", err)
	}

	server := Server{
		Addr:     addr,
		Listener: listener,
		IsClosed: &atomic.Bool{},
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

	err := response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Fatalf("unable to write status line, error: %v", err)
	}

	headers := response.GetDefaultHeaders(0)

	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatalf("unable to write headers, error: %v", err)
	}
}
