package server

import (
	"log"
	"log/slog"
	"net"
	"os"
)

var logger *slog.Logger = configLogger()

type HTTPServer struct {
	ListenAddr string
	Listener   net.Listener
	Ready chan struct{} // Signals server's readiness to accept connection
	Quit chan struct{} // Signals shutdown
}

// NewHTTPServer creates a new nexus server instance
// This functions accepts the address on which the server is to run
func NewHTTPServer(listenAddr string) *HTTPServer {
	return &HTTPServer{
		ListenAddr: listenAddr,
		Ready: make(chan struct{}),
		Quit: make(chan struct{}),
	}
}

// Starts an instance of the HTTP server
func (srv *HTTPServer) Start() error {

	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		logger.Error("unable to listen to connection", slog.Any("error", err))
		return err
	}

	srv.Listener = listener
	close(srv.Ready) // Signal readiness to handle connection

	srv.acceptConnection()

	<-srv.Quit // Block until a shutdown signal is received
	log.Println("Shutting down server...")
	srv.Listener.Close()

	return nil
}

func (srv *HTTPServer) acceptConnection() {
	log.Print("Ready to accept connection...")
	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			logger.Warn("could not accept connection", slog.Any("error", err))
			continue
		}

		go srv.handleConnection(conn)
	}
}

// Close signals that the server is shutting down
func (srv *HTTPServer) Close() {
	close(srv.Quit)
}

func (srv *HTTPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			res := "HTTP/1.1 500 Internal Server Error\r\n" +
			"Content-Type: text/plain\r\n" +
			"Connection: close\r\n\r\n" +
			"Error reading from connection\n"

			conn.Write([]byte(res))
		}

		res := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: close\r\n\r\n" +
		"Hello, World!\n"

		conn.Write([]byte(res))
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
