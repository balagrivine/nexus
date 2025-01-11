package main

import (
	"net"
	"fmt"
	"bufio"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	response, _ := bufio.NewReader(conn).ReadString('\n')

	fmt.Println(response)
}
