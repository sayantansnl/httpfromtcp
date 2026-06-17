package server

import (
	"bytes"
	"fmt"
	"github/sayantansnl/httpfromtcp/internal/request"
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatalf("unable to parse request from net.Conn, error: %v", err)
	}

	b := bytes.Buffer{}

	handlerError := s.Handler(&b, req)

	if handlerError.StatusCode == response.StatusServerError ||
		handlerError.StatusCode == response.StatusBadRequest {
		WriteError(conn, *handlerError)
		return
	}

	headers := response.GetDefaultHeaders(b.Len())

	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		log.Fatal("error in writing status line")
	}

	if err := response.WriteHeaders(conn, headers); err != nil {
		log.Fatal("error in writing headers")
	}

	conn.Write(b.Bytes())
}
