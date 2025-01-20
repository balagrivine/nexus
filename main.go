package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/balagrivine/nexus/server"
)

func main() {
	httpServer := server.NewHTTPServer("127.0.0.1:8080")

	if err := httpServer.Start(); err != nil {
		log.Fatal(err)
	}

	// Block indefinitely until server is ready to start
	// accepting connections
	<-httpServer.Ready

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan // Block until an interupt signal is received
	httpServer.Close()
}
