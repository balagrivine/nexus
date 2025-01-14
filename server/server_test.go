package server

import (
	"bufio"
	"net"
	"log"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Create a successful client connection", "405 Not Allowed\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				t.Errorf("unable to make a connection to the server %v\n", err)
			}
			if conn == nil {
				log.Fatalf("unable to make a connection to the server: connection refused")
			}
			defer conn.Close()

			rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
			if _, err := rw.WriteString("POST / HTTP/1.1\n"); err != nil {
				t.Errorf("unable to write to server, %v", err)
			}
			rw.Flush()

			response, err := rw.ReadString('\n')
			if err != nil {
				t.Errorf("error occured receiving response from server %v\n", err)
			}

			if response != test.want {
				t.Errorf("Expected %v got %v\n", test.want, response)
			}
		})
	}
}
