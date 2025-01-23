package http

import (
	"fmt"
	"net"
)

type HTTPStatus string

const (
	HTTP_200_OK           HTTPStatus = "200 OK"
	HTTP_201_CREATED      HTTPStatus = "201 CREATED"
	HTTP_404_NOT_FOUND    HTTPStatus = "404 NOT FOUND"
	HTTP_405_NOT_ALLOWED  HTTPStatus = "405 NOT ALLOWED"
	HTTP_400_BAD_REQUEST  HTTPStatus = "400 BAD REQUEST"
	HTTP_500_SERVER_ERROR HTTPStatus = "500 INTERNAL SERVER ERROR"
)

type Response struct {
	conn    net.Conn
	headers map[string]string
	status  HTTPStatus
}

// GetResponseWriter returns a pointer to Response struct
// used for composing HTTP responses
func GetResponseWriter(conn net.Conn) *Response {
	return &Response{
		conn: conn,
		headers: map[string]string {
			"Server": "nexus (Ubuntu)",
		},
		status: HTTP_200_OK,
	}
}

func (r *Response) setStatus(status HTTPStatus) {
	r.status = status
}

func (r *Response) AddHeader(key string, value string) {
	r.headers[key] = value
}

func (r *Response) write(content []byte) error {
	response := encode(content, r.status, r.headers)
	_, err := r.conn.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func (r *Response) Send(content []byte, status HTTPStatus) error {
	r.setStatus(status)
	if err := r.write(content); err != nil {
		return err
	}

	return nil
}

func encode(data []byte, status HTTPStatus, headers map[string]string) []byte {
	response := []byte{}

	response = append(response, VERSION...)
	response = append(response, SPACE...)
	response = append(response, []byte(fmt.Sprintf("%s\r\n", status))...)

	headers["Content-Length"] = fmt.Sprintf("%d", len(data))
	appendHeaders(&response, headers)

	response = append(response, CRLF...)
	response = append(response, data...)

	return response
}

func appendHeaders(response *[]byte, headers map[string]string) {
	for key, val := range headers {
		header := fmt.Sprintf("%s: %s\r\n", key, val)
		*response = append(*response, []byte(header)...)
	}
}
