// HTTP request reading and parsing

package http

import (
	"net"
	"strings"
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
}

// parseRequestLine parses "GET /foo HTTP/1.1" into its three parts.
func parseRequestLine(line string) (method, path, proto string, ok bool) {
	method, rest, ok1 := strings.Cut(line, " ")
	requestURI, proto, ok2 := strings.Cut(rest, " ")
	if !ok1 || !ok2 {
		return "", "", "", false
	}

	return method, requestURI, proto, true
}

func validMethod(method string) bool {
	return len(method) > 0
}

