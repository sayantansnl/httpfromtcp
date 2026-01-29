package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = ":42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if _, err := conn.Write([]byte(data)); err != nil {
			log.Fatalf("error in writing data to connection: %v", err)
		}
	}
}
