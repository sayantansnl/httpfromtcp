package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("\nerror in opening file: %v", err)
	}

	defer file.Close()

	fmt.Println("======Reading Data======")

	currentLine := ""

	for {
		buffer := make([]byte, 8)
		n, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("error in reading files: %v", err)
		}
		parts := strings.Split(string(buffer[:n]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", currentLine+parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
	if len(currentLine) > 0 {
		fmt.Printf("read: %s\n", currentLine)
	}
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	f.Read()
// }
