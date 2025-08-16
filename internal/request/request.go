package request

import (
	"bytes"
	"httpformscratch/internal/headers"
	"strconv"

	// "errors"
	"fmt"
	"io"

	// "strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
	Body          string
}

// func (r *RequestLine) validHTTP() bool {
// 	return r.HttpVersion == "HTTP/1.1"
// }

type Request struct {
	RequestLine RequestLine
	Headers *headers.Headers
	state parserState
	Body string
}

func getInt(headers *headers.Headers,name string, defaultValue int) int {
	valueStr,exists := headers.Get(name)
	if !exists {
		return defaultValue
	}

	value,err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
		Headers: headers.NewHeaders(),
		Body: "",
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
	stateHeaders parserState = "headers"
	StateBody parserState = "body"
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

func (r *Request) hasBody() bool {
	length := getInt(r.Headers,"content-length",0)
	return  length > 0 // should return true if body is other than 0 length
}

func (r *Request) parse(data []byte) (int, error){
	read := 0
outer: 
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break outer
		}
		switch r.state {
		case stateError:
			return 0, ErrRequestInErrorState
		case StateInit:
			rl,n,err := parseRequestLine(currentData)
			if err != nil {
				return 0,err
			}
			if n==0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = stateHeaders

			// r.state = StateDone -> why was i setting it done immediately ??????????????
		case stateHeaders:
			n,done,err := r.Headers.Parse(currentData)

			if err != nil {
				r.state = stateError
				return 0,err
			}

			if n == 0 {
				// return 0,nil
				break outer
			}
			read += n
			// in real world we would not get an eof after reading the data hence we can nicely transition to the body
			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
					// break
				}
			}
		case StateBody:
			length := getInt(r.Headers,"content-length",0)
			if length == 0 {
				// r.state = StateDone
				panic("chunked not implemented")
			}
			remaining := min(length - len(r.Body),len(currentData))
			r.Body += string(currentData[:remaining])

			if len(r.Body) == length {
				r.state =StateDone
			}
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
	for !request.done() {
        n, err := reader.Read(buf[bufLen:])
        if n > 0 {
            bufLen += n
            readN, err := request.parse(buf[:bufLen])
            if err != nil {
                return nil, err
            }
            copy(buf, buf[readN:bufLen])
            bufLen -= readN
        }
        if err != nil {
            if err == io.EOF {
                break // eof exit handled
            }
            return nil, err
        }
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