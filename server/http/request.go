// HTTP request reading and parsing

package http

import (
	"bytes"
	"errors"
	"io"
)

// A Request representsan HTTP request received
// by the server
type Request struct {

	// Method represents an HTTP method accompanying
	// the request
	Method string

	// Header contains the request header fields either
	// received by the server or to be sent by the client.
	Header map[string]string

	//Body is the request's body
	Body io.ReadCloser

	// The request path
	Path string

	// proto is the HTTP protocol used in the request
	Proto string
}

func Decode(content []byte) (*Request, error) {
	data := bytes.Split(content, CRLF)
	// TODO headers := make(map[string]string)

	requestLine := bytes.Split(data[0], SPACE)
	method := requestLine[0]
	path := requestLine[1]
	proto := requestLine[2]

	if !validMethod(method) {
		return nil, errors.New("invalid method")
	}

	return &Request{
		Method: string(method),
		Path:   string(path),
		Proto:  string(proto),
	}, nil

}

func validMethod(method []byte) bool {
	return string(method) == "GET"
}
