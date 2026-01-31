package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sayantansnl/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error in listening: %v", err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error in reading from request: %v", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		headers := req.Headers
		fmt.Println("Headers:")

		for fieldName, fieldVal := range headers {
			fmt.Printf("- %s: %s\n", fieldName, fieldVal)
		}

		fmt.Println("Body:")
		fmt.Printf("\n%s", string(req.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
