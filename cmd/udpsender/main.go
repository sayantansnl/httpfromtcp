package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const PORT = ":42069"
const NETWORK = "udp"

func main() {
	udpAddr, err := net.ResolveUDPAddr(NETWORK, PORT)
	if err != nil {
		log.Fatal("error in creating UDP Adress")
	}

	conn, err := net.DialUDP(NETWORK, nil, udpAddr)
	if err != nil {
		log.Fatal("error in establishing connection")
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error in reading: %v", err)
		}

		conn.Write([]byte(data))
	}
}
