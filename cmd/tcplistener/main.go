package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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

		fmt.Println("Connection accepted...")

		linesChannel := getLinesChannel(conn)
		for line := range linesChannel {
			fmt.Printf("%s\n", line)
		}

		fmt.Println("Connection closed...")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChannel := make(chan string)

	go func() {
		defer close(linesChannel)
		currentLine := ""

		for {
			bytes := make([]byte, 8)

			_, err := f.Read(bytes)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal("couldn't read said file")
				break
			}

			parts := strings.Split(string(bytes), "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLine += parts[i]
				linesChannel <- currentLine
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}

		if len(currentLine) > 0 {
			linesChannel <- currentLine
		}
	}()

	return linesChannel
}
