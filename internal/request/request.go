package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
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
}

var ErrMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrUnsupportedHTTPVersion = fmt.Errorf("unsuported http version")
var SEPERATOR = "\r\n"

func parseRequestLIne(b string) (*RequestLine,string,error){
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return  nil, b , nil
	}

	startLine := b[:idx]
	restOfMsg := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine," ")
	if len(parts) != 3 {
		return nil,restOfMsg, ErrMalformedRequestLine
	}

	httpParts := strings.Split(parts[2],"/")

	if len(httpParts) != 2  || httpParts[0] != "HTTP" || httpParts[1] != "1.1" {
		return nil,restOfMsg, ErrMalformedRequestLine
	}

	rl := &RequestLine{
		Method: parts[0],
		RequestTarget: parts[1],
		HttpVersion: httpParts[1],
	}

	// if !rl.validHTTP() {
	// 	return  nil,restOfMsg,ErrUnsupportedHTTPVersion
	// }

	return rl,restOfMsg,nil
}

func RequestFromReader(reader io.Reader) (*Request, error){
	data,err := io.ReadAll(reader)
	if err != nil {
		return nil,errors.Join(
			fmt.Errorf("unable to io.ReadAll"),
			err,
		)
	}
	str := string(data)
	rl, _, err := parseRequestLIne(str)
	if err != nil {
		return nil,err
	}

	return &Request{
		RequestLine: *rl,
	} ,err
}