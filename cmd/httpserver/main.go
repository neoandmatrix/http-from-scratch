package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpformscratch/internal/request"
	"httpformscratch/internal/response"
	"httpformscratch/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message: "Your problem not my problem\n",
			}
		case "/myproblem":
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message: "Woopise, my bad\n",
			}
		default:
			w.Write([]byte("All good frfr\n"))
		}
		return nil	
		
	})
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