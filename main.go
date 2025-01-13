package main

import (
	_ "os"

	"github.com/balagrivine/nexus/server"
)

func main() {

	srv := &server.TCPServer{}

	if err := srv.Run(); err != nil {
		srv.Close()
	}
	//if err := server.ListenAndServe(); err != nil {
	//	os.Exit(1)
	//}
}
