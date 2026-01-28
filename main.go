package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

	for fileSize > 0 {
		buffer := make([]byte, 8)
		numBytes, err := file.Read(buffer)
		if err != nil {
			log.Fatalf("\nerror in reading files: %v", err)
			break
		}
		fmt.Printf("read: %s\n", string(buffer))
		file.Truncate(8)
		fileSize -= int64(numBytes)
		if err == io.EOF {
			break
		}
	}
}
