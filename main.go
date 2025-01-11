package main

import (
	"net"
	"log/slog"
	"os"
)

func main() {
	if err := listenAndServe(); err != nil {
		os.Exit(1)
	}
}

func listenAndServe() error {

	const (
		Port = ":8080"
	)
	logger := configLogger()

	listener, err := net.Listen("tcp", Port)
	if err != nil {
		logger.Error("unable to listen to connection on port", slog.Any("error", err))
		return err
	}
	defer listener.Close()

	for {
		logger.Info("Accept a connection request")
		_, err := listener.Accept()
		if err != nil {
			logger.Warn("failed accepting a connection request", err)
			continue
		}
	}

	return nil
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
