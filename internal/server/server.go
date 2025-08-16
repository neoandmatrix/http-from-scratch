package server

import (
	"fmt"
	"httpformscratch/internal/response"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(_s *Server,conn io.ReadWriteCloser)  {
	defer conn.Close()
	// out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!`")
	// conn.Write(out)
	headers := response.GetDefaultHeaders(0)
	response.WriteStatusLine(conn,response.StatusOK)
	response.WriteHeaders(conn,headers)
	
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

func Serve(port uint16) (*Server,error){
	listner,err := net.Listen("tcp",fmt.Sprintf(":%d",port)) 
	if err != nil {
		return nil,err
	}
	server := &Server{closed: false}
	go runServer(server,listner)

	return server,nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}