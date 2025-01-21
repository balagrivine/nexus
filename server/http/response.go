package http

import (
	"fmt"
	"net"
)

type HTTPStatus string

const (
	HTTP_200_OK           HTTPStatus = "200"
	HTTP_201_CREATED      HTTPStatus = "201"
	HTTP_404_NOT_FOUND    HTTPStatus = "404"
	HTTP_405_NOT_ALLOWED  HTTPStatus = "405"
	HTTP_400_BAD_REQUEST  HTTPStatus = "400"
	HTTP_500_SERVER_ERROR HTTPStatus = "500"
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
		headers: make(map[string]string),
		status: HTTP_200_OK,
	}
}

func (r *Response) SetStatus(status HTTPStatus) {
	r.status = status
}

func (r *Response) AddHeader(key, value string) {
	r.headers[key] = value
}

func (r *Response) RemoveHeader(key string) {
	if _, present := r.headers[key]; !present {
		return
	}

	delete(r.headers, key)
}

func (r *Response) write(content []byte) {
	defer r.conn.Close()

	response := encode(content, r.status, r.headers)
	r.conn.Write(response)
}

func (r *Response) Send(content []byte, status HTTPStatus) {
	r.SetStatus(status)
	r.write(content)
}

func encode(data []byte, status HTTPStatus, headers map[string]string) []byte {
	response := []byte{}

	response = append(response, VERSION...)
	response = append(response, []byte(status)...)
	appendHeaders(&response, headers)

	response = append(response, CRLF...)
	response = append(response, CRLF...)
	response = append(response, data...)
	response = append(response, CRLF...)

	return response
}

func appendHeaders(response *[]byte, headers map[string]string) {
	for key, val := range headers {
		data := fmt.Sprintf("%s: %s", key, val)
		(*response) = append((*response), CRLF...)
		(*response) = append((*response), []byte(data)...)
	}
}
