package server

import (
	"bytes"
	"fmt"
	"httpformscratch/internal/request"
	"httpformscratch/internal/response"
	"io"
	"net"
)


type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer,req *request.Request ) *HandlerError


type Server struct {
	closed bool
	handler Handler
}

func runConnection(s *Server,conn io.ReadWriteCloser)  {
	defer conn.Close()
	// out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!`")
	// conn.Write(out)
	headers := response.GetDefaultHeaders(0)
	writer := bytes.NewBuffer([]byte{})
	r,err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn,response.StatusBadRequest)
		response.WriteHeaders(conn,headers)
		return
	}
	handlerError := s.handler(writer,r)

	var body []byte = nil
	var status response.StatusCode = response.StatusOK
	if handlerError != nil {
		body = []byte(handlerError.Message)
		response.WriteStatusLine(conn, handlerError.StatusCode)
		headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		response.WriteHeaders(conn, headers)
		conn.Write(body)
		return
} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length",fmt.Sprintf("%d",body))
	response.WriteStatusLine(conn,status)
	response.WriteHeaders(conn,headers)
	conn.Write(body)
	
}

func runServer(s *Server, listner net.Listener) {
		for {
			conn,err := listner.Accept()
			if s.closed {
				return
			}
			if err != nil {
				return
			}
		go runConnection(s,conn)
		}
}

func Serve(port uint16, handler Handler) (*Server,error){
	listner,err := net.Listen("tcp",fmt.Sprintf(":%d",port)) 
	if err != nil {
		return nil,err
	}
	server := &Server{closed: false,handler: handler}
	go runServer(server,listner)

	return server,nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}