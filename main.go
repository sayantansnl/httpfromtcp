package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Printf("\nerror in opening file: %v", err)
	}

	fmt.Println("======Reading Data======")

	channel := getLinesChannel(file)

	for line := range channel {
		fmt.Printf("read: %s\n", line)
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

	for {
		buffer := make([]byte, 8)
		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Errorf("error in reading files, %w", err)
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
	close(ch)
}
