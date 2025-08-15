package request

import (
	"bytes"
	// "errors"
	"fmt"
	"io"
	// "strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// func (r *RequestLine) validHTTP() bool {
// 	return r.HttpVersion == "HTTP/1.1"
// }

type Request struct {
	RequestLine RequestLine
	state parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHTTPVersion = fmt.Errorf("unsuported http version")
var ErrRequestInErrorState = fmt.Errorf("request in error state")
var SEPERATOR = []byte("\r\n")

type parserState string
const (
	StateInit parserState = "init"
	StateDone parserState = "done"
	stateError parserState = "error"
)

func parseRequestLine(b []byte) (*RequestLine,int,error){
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return  nil, 0 , nil
	}

	startLine := b[:idx]
	read := idx+len(SEPERATOR)

	parts := bytes.Split(startLine,[]byte(" "))
	if len(parts) != 3 {
		return nil,0, ErrMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2],[]byte("/"))

	if len(httpParts) != 2  || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil,0, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method: 		string(parts[0]),
		RequestTarget: 	string(parts[1]),
		HttpVersion: 	string(httpParts[1]),
	}

	// if !rl.validHTTP() {
	// 	return  nil,restOfMsg,ErrUnsupportedHTTPVersion
	// }

	return rl,read,nil
}

func (r *Request) parse(data []byte) (int, error){
	read := 0
outer: 
	for {
		switch r.state {
		case stateError:
			return 0, ErrRequestInErrorState
		case StateInit:
			rl,n,err := parseRequestLine(data[read:])
			if err != nil {
				return 0,err
			}
			if n==0 {
				break outer
			}
			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
		}
	} 
	return read,nil
}

func (r *Request) done() bool{
	return r.state == StateDone
}

// func (r *Request) error() bool{
// 	return r.state == stateError || r.state == stateError
// }

func RequestFromReader(reader io.Reader) (*Request, error){
	request := newRequest()
	// buffer could get over run
	buf := make([]byte,1024)
	// var request *Request
	bufLen := 0
	// simulating slowly reading one at a time
	for  !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil{
			return nil,err
		}
		bufLen += n
		readN,err := request.parse(buf[:bufLen])
		if err != nil {
			return nil,err
		}
		// data,err := io.ReadAll(reader)
		
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	
	}

	// if err != nil {
	// 	return nil,errors.Join(
	// 		fmt.Errorf("unable to io.ReadAll"),
	// 		err,
	// 	)

	// 	buf
	// }
	// str := string(data)
	// rl, _, err := parseRequestLine(str)
	// if err != nil {
	// 	return nil,err
	// }

	// return &Request{
	// 	RequestLine: *rl,
	// } ,err

	return  request,nil

}