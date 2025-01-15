package main

import (
	"log"

	"github.com/balagrivine/nexus/server"
)

func main() {

	httpServer := server.NewHTTPServer("127.0.0.1:8080")

	if err := httpServer.Start(); err != nil {
		httpServer.Close()
		log.Fatal(err)
	}
}
