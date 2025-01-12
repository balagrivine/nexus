package main

import (
	"os"

	"github.com/balagrivine/server"
)

func main() {
	if err := server.ListenAndServe(); err != nil {
		os.Exit(1)
	}
}
