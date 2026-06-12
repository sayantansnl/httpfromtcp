package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("couldn't open said file")
	}

	defer file.Close()

	linesChannel := getLinesChannel(file)

	for thing := range linesChannel {
		fmt.Printf("read: %s\n", thing)
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
