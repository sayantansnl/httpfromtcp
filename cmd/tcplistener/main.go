package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		channel := getLinesChannel(conn)

		for line := range channel {
			fmt.Println(line)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	channel := make(chan string)
	go receiveData(channel, f)
	return channel
}

func receiveData(ch chan string, f io.ReadCloser) {
	currentLine := ""
	defer f.Close()
	defer close(ch)

	for {
		buffer := make([]byte, 8)
		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("error in reading files, %v", err)
			return
		}
		parts := strings.Split(string(buffer[:n]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			ch <- currentLine + parts[i]
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
	if len(currentLine) > 0 {
		ch <- currentLine
	}
}
