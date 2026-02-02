package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sayantansnl/httpfromtcp/internal/request"
	"github.com/sayantansnl/httpfromtcp/internal/response"
	"github.com/sayantansnl/httpfromtcp/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	reqTarget := req.RequestLine.RequestTarget

	if reqTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeBadRequest,
			Message:    "Your problem is not my problem\r\n",
		}
	}

	if reqTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusCodeServerError,
			Message:    "Woopsie, my bad\r\n",
		}
	}

	w.Write([]byte("All good, frfr\r\n"))
	return nil
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
