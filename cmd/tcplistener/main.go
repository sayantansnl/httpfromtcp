package main

import (
	"fmt"
	"log"
	"net"

	"github/sayantansnl/httpfromtcp/internal/request"
)

const PORT = ":42069"
const NETWORK = "tcp"

func main() {
	listener, err := net.Listen(NETWORK, PORT)
	if err != nil {
		log.Fatalf("error in creating listener, error: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error in creating a connection, error: %v", err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error in parsing request from main, error: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s'\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s'\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s'\n", req.RequestLine.HttpVersion)

		fmt.Println("Headers:")

		for key, val := range req.Headers {
			fmt.Printf("- %s: %s\n", key, val)
		}
	}
}
