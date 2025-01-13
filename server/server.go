package server

import (
	"bufio"
	"strings"
	"log/slog"
	"net"
	"os"
)

var logger *slog.Logger = configLogger()

type Server interface {
	Run() error
	Close() error
}

type TCPServer struct {
	port string
	server net.Listener

}

// Run starts the tcp server
func (t *TCPServer) Run() (err error) {

	t.port = ":8080"

	t.server, err = net.Listen("tcp", t.port)
	if err != nil {
		logger.Error("unable to listen to connection on port", slog.Any("error", err))
		return
	}

	for {
		logger.Info("Accept a connection request.")
		conn, err := t.server.Accept()
		if err != nil {
			logger.Warn("failed accepting a connection request", err)
			continue
		}

		logger.Info("Handle incoming messages.")
		go t.handleConnection(conn)
	}

	return
}

func (t *TCPServer) Close() {
	// Check for nil value to avoid nil pointer
	// dereference which causes the program to panic
	if t.server != nil {
		t.server.Close()
	}
	os.Exit(1)
}

func (t *TCPServer) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	req, err := rw.ReadString('\n')
	if err != nil {
		logger.Error("error reading request from client", slog.Any("error", err))
	}

	// Split the string requests to obtain
	// HTTP Method and requested path
	parts := strings.Split(req, " ")

	if _, err = rw.WriteString("HTTP/1.1 200 OK\n"); err != nil {
		logger.Warn("cannot write to connection", slog.Any("error", err))
	}

	if err = rw.Flush(); err != nil {
		logger.Warn("failed flush", slog.Any("error", err))
	}
}

// Creates and returns a custom Logger instance
func configLogger() *slog.Logger {
	logLevel := &slog.LevelVar{}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger
}
