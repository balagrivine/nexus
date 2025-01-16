package server

import (
	"testing"
	"net"
	"bufio"
	"fmt"
)

func TestHTTPServer(t *testing.T) {
	tests := []struct{
		name string
		input string
		want string
	}{
		{"Successfully connect to the server", "\n", "HTTP/1.1 200 OK\n"},
		{"Successfully ping the server", "PING\n", "PONG\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, _ := net.Dial("tcp", "127.0.0.1:8080")
			defer conn.Close()

			fmt.Fprintf(conn, test.input)

			res, _ := bufio.NewReader(conn).ReadString('\n')
			if test.want != res {
				t.Errorf("Expected %v, got %v", test.want, res)
			}
		})
	}
}
