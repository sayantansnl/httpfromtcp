package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sayantansnl/httpfromtcp/internal/headers"
	"github.com/sayantansnl/httpfromtcp/internal/request"
	"github.com/sayantansnl/httpfromtcp/internal/response"
	"github.com/sayantansnl/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	reqTarget := req.RequestLine.RequestTarget

	if reqTarget == "/yourproblem" {
		handleBadRequest(w, req)
		return
	}

	if reqTarget == "/myproblem" {
		handleServerError(w, req)
		return
	}

	if strings.HasPrefix(reqTarget, "/httpbin") {
		handleHTTPBinProxy(w, req)
		return
	}

	handleSuccess(w, req)
}

func handleBadRequest(w *response.Writer, _ *request.Request) {
	respBody := "<html><head><title>400 Bad Request</title></head><body><h1>Bad Request</h1><p>Your request honestly kinda sucked.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeBadRequest)
	headers := response.GetDefaultHeaders(len([]byte(respBody)))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(respBody))
}

func handleServerError(w *response.Writer, _ *request.Request) {
	respBody := "<html><head><title>500 Internal Server Error</title></head><body><h1>Internal Server Error</h1><p>Okay, you know what? This one is on me.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeServerError)
	headers := response.GetDefaultHeaders(len([]byte(respBody)))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(respBody))
}

func handleSuccess(w *response.Writer, _ *request.Request) {
	respBody := "<html><head><title>200 OK</title></head><body><h1>Success!</h1><p>Your request was an absolute banger.</p></body></html>"
	w.WriteStatusLine(response.StatusCodeSuccess)
	headers := response.GetDefaultHeaders(len([]byte(respBody)))
	headers.Override("Content-Type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(respBody))
}

func handleHTTPBinProxy(w *response.Writer, req *request.Request) {
	trimmed := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	newUrl := fmt.Sprintf("https://httpbin.org%s", trimmed)

	res, err := http.Get(newUrl)
	if err != nil {
		log.Printf("proxy error: unable to get response from %s: %v", newUrl, err)
		w.WriteStatusLine(response.StatusCodeServerError)
		body := []byte("upstream error")
		h := response.GetDefaultHeaders(len(body))
		w.WriteHeaders(h)
		w.WriteBody(body)
		return
	}

	defer res.Body.Close()

	w.WriteStatusLine(response.StatusCodeSuccess)
	respHeaders := headers.NewHeaders()
	respHeaders.Set("Transfer-Encoding", "chunked")
	respHeaders.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteHeaders(respHeaders)

	buff := make([]byte, 1024)
	fullBody := []byte{}
	for {
		n, err := res.Body.Read(buff)
		if n > 0 {
			fullBody = append(fullBody, buff[:n]...)
			_, writeErr := w.WriteChunkedBody(buff[:n])
			if writeErr != nil {
				log.Printf("unable to write to buffer: %v", writeErr)
				return
			}
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("unable to write to buffer: %v", err)
			return
		}
	}
	w.WriteChunkedBodyDone()
	sum := sha256.Sum256(fullBody)

	trailerHeaders := headers.NewHeaders()
	trailerHeaders.Set("X-Content-SHA256", hex.EncodeToString(sum[:]))
	trailerHeaders.Set("X-Content-Length", fmt.Sprint(len(fullBody)))

	w.WriteTrailers(trailerHeaders)
}
