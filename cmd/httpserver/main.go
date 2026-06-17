package main

import (
	"github/sayantansnl/httpfromtcp/internal/request"
	"github/sayantansnl/httpfromtcp/internal/response"
	"github/sayantansnl/httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerFunc)
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

func handlerFunc(w io.Writer, req *request.Request) *server.HandlerError {
	handlerError := server.HandlerError{}

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handlerError = server.HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    "Your problem is not my problem\n",
		}

	case "/myproblem":
		handlerError = server.HandlerError{
			StatusCode: response.StatusServerError,
			Message:    "Woopsie, my bad\n",
		}

	default:
		handlerError = server.HandlerError{
			StatusCode: response.StatusOK,
			Message:    "All good, frfr\n",
		}
	}

	w.Write([]byte(handlerError.Message))
	return &handlerError
}
