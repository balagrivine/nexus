package main

import (
	"os"
)

func main() {
	if err := listenAndServe(); err != nil {
		os.Exit(1)
	}
}
