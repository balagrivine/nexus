package server

import (
	"bufio"
	"log/slog"
	"net"
	"os"
)

var logger *slog.Logger = configLogger()

func ListenAndServe() error {

	const (
		Port = ":8080"
	)

	listener, err := net.Listen("tcp", Port)
	if err != nil {
		logger.Error("unable to listen to connection on port", slog.Any("error", err))
		return err
	}
	defer listener.Close()

	for {
		logger.Info("Accept a connection request.")
		conn, err := listener.Accept()
		if err != nil {
			logger.Warn("failed accepting a connection request", err)
			continue
		}

		logger.Info("Handle incoming messages.")
		go handleConnection(conn)
	}

	return nil
}

func handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	_, err := rw.ReadString('\n')
	if err != nil {
		logger.Error("error reading request from client", slog.Any("error", err))
	}

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
