package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Connect to the server", "GET / HTTP/1.1\r\n", "HTTP/1.1 200 OK\r\n"},
	}

	mockSrv := NewHTTPServer("127.0.0.1:8080")

	// Run the mock server in its goroutine to prevent blocking
	go func() {
		if err := mockSrv.Start(); err != nil {
			t.Errorf("Error starting server: %v\n", err)
			return
		}
	}()

	// Wait for server to signal readiness
	<-mockSrv.Ready

	t.Cleanup(func() {
		mockSrv.Close()
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", "127.0.0.1:8080")
			if err != nil {
				t.Errorf(fmt.Sprintf("%v", err))
				return
			}

			fmt.Fprintf(conn, test.input)

			res, _ := bufio.NewReader(conn).ReadString('\n')
			if test.want != res {
				t.Errorf("Expected %v, got %v", test.want, res)
			}
			conn.Close()
		})
	}
}
