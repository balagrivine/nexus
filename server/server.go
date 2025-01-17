package server

import (
	"bufio"
	"log"
	"log/slog"
	"net"
	"os"
)

var logger *slog.Logger = configLogger()

type HTTPServer struct {
	ListenAddr string
	Listener   net.Listener
}

// NewHTTPServer creates a new nexus server instance
// This functions accepts the address on which the server is to run
func NewHTTPServer(listenAddr string) *HTTPServer {
	return &HTTPServer{
		ListenAddr: listenAddr,
	}
}

// Starts an instance of the HTTP server
func (srv *HTTPServer) Start() error {

	listener, err := net.Listen("tcp", srv.ListenAddr)
	if err != nil {
		logger.Error("unable to listen to connection on port", slog.Any("error", err))
		return err
	}

	srv.Listener = listener

	go srv.acceptConnection()
	return nil
}

func (srv *HTTPServer) acceptConnection() {
	log.Print("Ready to accept connection")
	for {
		conn, err := srv.Listener.Accept()
		if err != nil {
			logger.Warn("could not accept connection", slog.Any("error", err))
			return
		}

		go srv.handleConnection(conn)
	}
	return
}

// Shuts down the current running instance of the server
func (srv *HTTPServer) Close() {
	srv.Listener.Close()
}

func (srv *HTTPServer) handleConnection(conn net.Conn) error {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				logger.Info("Connection reset by client")
				break
			}
			logger.Error("error reading request from client", slog.Any("error", err))
			return err
		}

		if req == "PING\n" {
			srv.handlePing(rw)
			break
		}

		if _, err = rw.WriteString("HTTP/1.1 200 OK\n"); err != nil {
			logger.Warn("cannot write to connection", slog.Any("error", err))
		}

		if err = rw.Flush(); err != nil {
			logger.Warn("failed flush", slog.Any("error", err))
		}
	}
	return nil
}

// Responds to a PING request with a PONG
func (srv *HTTPServer) handlePing(rw *bufio.ReadWriter) {
	if _, err := rw.WriteString("PONG\n"); err != nil {
		logger.Warn("cannot write to connection", "error", err)
	}
	rw.Flush()
	return
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
