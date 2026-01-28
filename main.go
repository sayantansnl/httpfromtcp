package main

import (
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

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("\nerror getting file info: %v", err)
	}

	fileSize := fileInfo.Size()
	defer file.Close()

	fmt.Println("======Reading Data======")

	currentLine := ""

	for fileSize > 0 {
		buffer := make([]byte, 8)
		numBytes, err := file.Read(buffer)
		if err != nil {
			log.Fatalf("\nerror in reading files: %v", err)
			break
		}
		parts := strings.Split(string(buffer), "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s\n", currentLine+parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
		fileSize -= int64(numBytes)
		if err == io.EOF {
			break
		}
	}
	if len(currentLine) > 0 {
		fmt.Printf("read: %s\n", currentLine)
	}
}
